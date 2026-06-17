package dto

type CreateSettlementRequest struct {
	Amount      int64  `json:"amount"`
	BankAccount string `json:"bank_account"`
}
