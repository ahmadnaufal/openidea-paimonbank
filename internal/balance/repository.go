package balance

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type balanceRepo struct {
	db *sqlx.DB
}

func NewBalanceRepo(db *sqlx.DB) balanceRepo {
	return balanceRepo{db: db}
}

func (r *balanceRepo) AddBalance(ctx context.Context, tx *sql.Tx, val BalanceHistory) error {
	baseQuery := `
		INSERT INTO
			balance_histories
			(id, user_id, currency, balance, source_bank_account_number, source_bank_name, transfer_proof_img_url)
		VALUES
			(:id, :user_id, :currency, :balance, :source_bank_account_number, :source_bank_name, :transfer_proof_img_url)
	`

	query, args, err := sqlx.Named(baseQuery, val)
	if err != nil {
		return err
	}

	dollarQuery := sqlx.Rebind(sqlx.DOLLAR, query)

	if tx != nil {
		_, err = tx.ExecContext(ctx, dollarQuery, args...)
	} else {
		_, err = r.db.ExecContext(ctx, dollarQuery, args...)
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *balanceRepo) GetBalanceHistory(ctx context.Context, payload GetBalanceHistoryRequest) ([]BalanceHistory, uint, error) {
	var balanceHistories []BalanceHistory

	baseQuery := `
		SELECT
			bh.id,
			bh.user_id,
			bh.currency,
			bh.balance,
			bh.source_bank_account_number,
			bh.source_bank_name,
			bh.transfer_proof_img_url,
			bh.created_at
		FROM
			balance_histories bh
		WHERE
			bh.user_id = ?
		%s
	`

	args := []interface{}{payload.UserID}

	filterQuery, filterArgs := getFilter(payload)

	args = append(args, filterArgs...)

	queryWithFilter := fmt.Sprintf(baseQuery, filterQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS temp", queryWithFilter)

	var count uint
	err := r.db.GetContext(ctx, &count, sqlx.Rebind(sqlx.DOLLAR, countQuery), args...)
	if err != nil {
		return balanceHistories, count, err
	}

	orderQuery := getSortBy(payload)
	limitQuery, limitArgs := getLimitAndOffset(payload)
	args = append(args, limitArgs...)

	query := fmt.Sprintf("%s %s %s", queryWithFilter, orderQuery, limitQuery)

	err = r.db.SelectContext(ctx, &balanceHistories, sqlx.Rebind(sqlx.DOLLAR, query), args...)
	if err != nil {
		return balanceHistories, count, err
	}

	return balanceHistories, count, nil
}

func getFilter(_ GetBalanceHistoryRequest) (string, []interface{}) {
	args := []interface{}{}
	filter := ""
	return filter, args
}

func getSortBy(_ GetBalanceHistoryRequest) string {
	// hardcoded for now
	return `ORDER BY bh.created_at DESC`
}

func getLimitAndOffset(req GetBalanceHistoryRequest) (string, []interface{}) {
	// by default, set limit to 50
	query := "LIMIT ? OFFSET ?"

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// offset by default will be 0
	offset := req.Offset

	args := []interface{}{limit, offset}

	return query, args
}

func (r *balanceRepo) GetBalancePerCurrencies(ctx context.Context, userID string) ([]BalancePerCurrency, error) {
	var balancePerCurrency []BalancePerCurrency

	baseQuery := `
		SELECT
			currency,
			COALESCE(SUM(balance), 0) AS balance_per_currency
		FROM
			balance_histories
		WHERE
			user_id = $1
		GROUP BY
			currency
		ORDER BY
			balance_per_currency DESC
	`

	err := r.db.SelectContext(ctx, &balancePerCurrency, baseQuery, userID)
	if err != nil {
		return balancePerCurrency, err
	}

	return balancePerCurrency, nil
}
