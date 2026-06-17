package entity

import "time"

type AuditLog struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ActorID    string    `gorm:"type:uuid;index" json:"actor_id"`
	MerchantID string    `gorm:"type:uuid;index" json:"merchant_id"`
	Action     string    `gorm:"not null;index" json:"action"`
	Resource   string    `gorm:"not null" json:"resource"`
	Metadata   string    `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
