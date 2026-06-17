package entity

import "time"

type Merchant struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     string    `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Name       string    `gorm:"not null" json:"name"`
	Email      string    `gorm:"not null" json:"email"`
	WebhookURL string    `json:"webhook_url"`
	SecretKey  string    `gorm:"not null" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Merchant) TableName() string { return "merchants" }

type APIKey struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MerchantID string    `gorm:"type:uuid;not null;index" json:"merchant_id"`
	Key        string    `gorm:"uniqueIndex;not null" json:"key"`
	Active     bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

func (APIKey) TableName() string { return "api_keys" }
