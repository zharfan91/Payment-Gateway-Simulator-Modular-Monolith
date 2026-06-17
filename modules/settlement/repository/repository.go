package repository

import (
	"context"

	"github.com/zharf/payment-gateway-simulator/modules/settlement/entity"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, settlement *entity.Settlement) error {
	return r.db.WithContext(ctx).Create(settlement).Error
}

func (r *Repository) ListByMerchant(ctx context.Context, merchantID string) ([]entity.Settlement, error) {
	var settlements []entity.Settlement
	err := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).Order("created_at DESC").Find(&settlements).Error
	return settlements, err
}
