package entity

import "time"

type WebhookLog struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID string    `gorm:"type:uuid;not null;index" json:"merchant_id"`
	Event      string    `gorm:"not null" json:"event"`
	URL        string    `gorm:"not null" json:"url"`
	Payload    string    `gorm:"type:jsonb;not null" json:"payload"`
	StatusCode int       `json:"status_code"`
	Attempt    int       `gorm:"not null" json:"attempt"`
	Success    bool      `gorm:"not null" json:"success"`
	Error      string    `json:"error"`
	CreatedAt  time.Time `json:"created_at"`
}

func (WebhookLog) TableName() string { return "webhook_logs" }
