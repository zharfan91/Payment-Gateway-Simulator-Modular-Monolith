package repository

import (
	"context"

	"github.com/zharf/payment-gateway-simulator/modules/payment/entity"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) DB() *gorm.DB { return r.db }

func (r *Repository) Create(ctx context.Context, payment *entity.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *Repository) Update(ctx context.Context, payment *entity.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *Repository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error
	return &payment, err
}

func (r *Repository) FindByIdempotencyKey(ctx context.Context, merchantID, key string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithContext(ctx).Where("merchant_id = ? AND idempotency_key = ?", merchantID, key).First(&payment).Error
	return &payment, err
}

func (r *Repository) CreateLog(ctx context.Context, log *entity.PaymentLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
