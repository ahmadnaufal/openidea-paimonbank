package balance

import "github.com/jmoiron/sqlx"

type balanceRepo struct {
	db *sqlx.DB
}

func NewBalanceRepo(db *sqlx.DB) balanceRepo {
	return balanceRepo{db: db}
}
