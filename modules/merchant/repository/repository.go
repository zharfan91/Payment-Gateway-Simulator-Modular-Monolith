package repository

import (
	"context"

	"github.com/zharf/payment-gateway-simulator/modules/merchant/entity"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, merchant *entity.Merchant) error {
	return r.db.WithContext(ctx).Create(merchant).Error
}

func (r *Repository) Update(ctx context.Context, merchant *entity.Merchant) error {
	return r.db.WithContext(ctx).Save(merchant).Error
}

func (r *Repository) FindByUserID(ctx context.Context, userID string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&merchant).Error
	return &merchant, err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.db.WithContext(ctx).First(&merchant, "id = ?", id).Error
	return &merchant, err
}

func (r *Repository) CreateAPIKey(ctx context.Context, key *entity.APIKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}

func (r *Repository) FindAPIKey(ctx context.Context, key string) (*entity.APIKey, error) {
	var apiKey entity.APIKey
	err := r.db.WithContext(ctx).Where("key = ? AND active = true", key).First(&apiKey).Error
	return &apiKey, err
}
