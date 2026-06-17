package repository

import (
	"context"

	"github.com/zharf/payment-gateway-simulator/modules/notification/entity"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreateWebhookLog(ctx context.Context, log *entity.WebhookLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
