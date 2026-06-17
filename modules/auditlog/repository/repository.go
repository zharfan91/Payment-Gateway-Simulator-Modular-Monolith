package repository

import (
	"context"

	"github.com/zharf/payment-gateway-simulator/modules/auditlog/entity"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, log *entity.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *Repository) ListByMerchant(ctx context.Context, merchantID string) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}
