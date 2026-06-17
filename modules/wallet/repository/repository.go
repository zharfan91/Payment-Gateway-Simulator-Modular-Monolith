package repository

import (
	"context"

	"github.com/zharf/payment-gateway-simulator/modules/wallet/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) DB() *gorm.DB { return r.db }

func (r *Repository) Create(ctx context.Context, wallet *entity.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *Repository) FindByMerchantID(ctx context.Context, merchantID string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	err := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID).First(&wallet).Error
	return &wallet, err
}

func (r *Repository) FindByMerchantIDForUpdate(ctx context.Context, tx *gorm.DB, merchantID string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("merchant_id = ?", merchantID).First(&wallet).Error
	return &wallet, err
}

func (r *Repository) Save(ctx context.Context, tx *gorm.DB, wallet *entity.Wallet) error {
	return tx.WithContext(ctx).Save(wallet).Error
}

func (r *Repository) CreateTransaction(ctx context.Context, tx *gorm.DB, trx *entity.WalletTransaction) error {
	return tx.WithContext(ctx).Create(trx).Error
}
