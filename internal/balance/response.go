package balance

type BalanceSourceResponse struct {
	BankAccountNumber string `json:"bankAccountNumber"`
	BankName          string `json:"bankName"`
}

type BalanceHistoryResponse struct {
	TransactionID    string                `json:"transactionId"`
	Balance          int                   `json:"balance"`
	Currency         string                `json:"currency"`
	TransferProofImg string                `json:"transferProofImg"`
	CreatedAt        uint64                `json:"createdAt"`
	Source           BalanceSourceResponse `json:"source"`
}

type CurrencyBalanceResponse struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}
