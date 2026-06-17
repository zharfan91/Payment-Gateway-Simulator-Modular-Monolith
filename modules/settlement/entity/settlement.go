package entity

import "time"

type Status string

const (
	StatusPending Status = "PENDING"
	StatusPaid    Status = "PAID"
	StatusFailed  Status = "FAILED"
)

type Settlement struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID  string    `gorm:"type:uuid;not null;index" json:"merchant_id"`
	Amount      int64     `gorm:"not null" json:"amount"`
	BankAccount string    `gorm:"not null" json:"bank_account"`
	Status      Status    `gorm:"not null" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Settlement) TableName() string { return "settlements" }
