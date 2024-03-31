package balance

import (
	"errors"
	"time"
)

type AddBalanceRequest struct {
	SenderBankAccountNumber string `json:"senderBankAccountNumber" validate:"required,min=5,max=30"`
	SenderBankName          string `json:"senderBankName" validate:"required,min=5,max=30"`
	AddedBalance            uint   `json:"addedBalance" validate:"required,gte=0"`
	Currency                string `json:"currency" validate:"required,iso4217"`
	TransferProofImg        string `json:"transferProofImg" validate:"required,url"`

	UserID string
}

type GetBalanceHistoryRequest struct {
	Limit  uint `query:"limit"`
	Offset uint `query:"offset"`

	UserID  string
	Queries map[string]string
}

// Validate is a function for additional validation related to query
func (r *GetBalanceHistoryRequest) Validate() error {
	queries := r.Queries

	if val, ok := queries["limit"]; ok && val == "" {
		return errors.New("limit is empty")
	}

	if val, ok := queries["offset"]; ok && val == "" {
		return errors.New("offset is empty")
	}

	return nil
}

type CreateTransactionRequest struct {
	RecipientBankAccountNumber string `json:"recipientBankAccountNumber" validate:"required,min=5,max=30"`
	RecipientBankName          string `json:"recipientBankName" validate:"required,min=5,max=30"`
	FromCurrency               string `json:"fromCurrency" validate:"required,iso4217"`
	Balances                   uint   `json:"balances" validate:"required,gt=0"`

	UserID string
}

type BalanceHistory struct {
	ID                      string    `db:"id"`
	UserID                  string    `db:"user_id"`
	Currency                string    `db:"currency"`
	Balance                 int       `db:"balance"`
	SourceBankAccountNumber string    `db:"source_bank_account_number"`
	SourceBankName          string    `db:"source_bank_name"`
	TransferProofImg        string    `db:"transfer_proof_img_url"`
	CreatedAt               time.Time `db:"created_at"`
}

type BalancePerCurrency struct {
	Balance  int    `db:"balance_per_currency"`
	Currency string `db:"currency"`
}
