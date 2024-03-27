package balance

type AddBalanceRequest struct {
	SenderBankAccountNumber string `json:"senderBankAccountNumber" validate:"required,min=5,max=30"`
	SenderBankName          string `json:"senderBankName" validate:"required,min=5,max=30"`
	AddedBalance            uint   `json:"addedBalance" validate:"required,gte=0"`
	Currency                string `json:"currency" validate:"required,currency"`
	TransferProofImg        string `json:"transferProofImg" validate:"required,url"`
}
