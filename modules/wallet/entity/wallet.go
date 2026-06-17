package entity

import "time"

type Wallet struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID string    `gorm:"type:uuid;not null;uniqueIndex" json:"merchant_id"`
	Balance    int64     `gorm:"not null;default:0" json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Wallet) TableName() string { return "wallets" }

type TransactionType string

const (
	TxCredit     TransactionType = "CREDIT"
	TxDebit      TransactionType = "DEBIT"
	TxAdjustment TransactionType = "ADJUSTMENT"
)

type WalletTransaction struct {
	ID           string          `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	WalletID     string          `gorm:"type:uuid;not null;index" json:"wallet_id"`
	PaymentID    string          `gorm:"type:uuid;index" json:"payment_id"`
	SettlementID string          `gorm:"type:uuid;index" json:"settlement_id"`
	Type         TransactionType `gorm:"not null" json:"type"`
	Amount       int64           `gorm:"not null" json:"amount"`
	Balance      int64           `gorm:"not null" json:"balance"`
	Description  string          `json:"description"`
	CreatedAt    time.Time       `json:"created_at"`
}

func (WalletTransaction) TableName() string { return "wallet_transactions" }
