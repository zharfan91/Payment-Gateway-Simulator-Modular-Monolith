package entity

import "time"

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusSuccess  Status = "SUCCESS"
	StatusFailed   Status = "FAILED"
	StatusExpired  Status = "EXPIRED"
	StatusRefunded Status = "REFUNDED"
)

type Payment struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID     string    `gorm:"type:uuid;not null;index" json:"merchant_id"`
	ExternalID     string    `gorm:"index" json:"external_id"`
	Amount         int64     `gorm:"not null" json:"amount"`
	Fee            int64     `gorm:"not null" json:"fee"`
	NetAmount      int64     `gorm:"not null" json:"net_amount"`
	Status         Status    `gorm:"not null;index" json:"status"`
	PaymentToken   string    `gorm:"uniqueIndex;not null" json:"payment_token"`
	PaymentURL     string    `gorm:"not null" json:"payment_url"`
	Description    string    `json:"description"`
	IdempotencyKey string    `gorm:"index" json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Payment) TableName() string { return "payments" }

type PaymentLog struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PaymentID string    `gorm:"type:uuid;not null;index" json:"payment_id"`
	Status    Status    `gorm:"not null" json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (PaymentLog) TableName() string { return "payment_logs" }
