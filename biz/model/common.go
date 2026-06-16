package model

type Account struct {
	AccountID int64  `json:"account_id"`
	Balance   string `json:"balance"`
}

type EmptyResponse struct{}

type GetAccount struct {
	AccountID int64 `json:"account_id"`
}

type Transaction struct {
	TransactionID        string `json:"transaction_id,omitempty"`
	SourceAccountID      int64  `json:"source_account_id"`
	DestinationAccountID int64  `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

type NewAccount struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}
