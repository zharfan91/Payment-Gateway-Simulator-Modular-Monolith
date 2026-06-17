package service

import (
	"context"
	"errors"

	"github.com/zharf/payment-gateway-simulator/modules/wallet/entity"
	"github.com/zharf/payment-gateway-simulator/modules/wallet/repository"
	"gorm.io/gorm"
)

const FeeBasisPoints int64 = 290

var ErrInsufficientBalance = errors.New("insufficient balance")

type Service struct{ repo *repository.Repository }

func New(repo *repository.Repository) *Service { return &Service{repo: repo} }

func CalculateFee(amount int64) (fee int64, net int64) {
	fee = amount * FeeBasisPoints / 10000
	return fee, amount - fee
}

func (s *Service) Ensure(ctx context.Context, merchantID string) error {
	_, err := s.repo.FindByMerchantID(ctx, merchantID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return s.repo.Create(ctx, &entity.Wallet{MerchantID: merchantID})
}

func (s *Service) Get(ctx context.Context, merchantID string) (*entity.Wallet, error) {
	return s.repo.FindByMerchantID(ctx, merchantID)
}

func (s *Service) Credit(ctx context.Context, merchantID, paymentID string, amount int64) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		wallet, err := s.repo.FindByMerchantIDForUpdate(ctx, tx, merchantID)
		if err != nil {
			return err
		}
		wallet.Balance += amount
		if err := s.repo.Save(ctx, tx, wallet); err != nil {
			return err
		}
		return s.repo.CreateTransaction(ctx, tx, &entity.WalletTransaction{
			WalletID: wallet.ID, PaymentID: paymentID, Type: entity.TxCredit, Amount: amount, Balance: wallet.Balance, Description: "payment captured",
		})
	})
}

func (s *Service) Debit(ctx context.Context, merchantID, settlementID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		wallet, err := s.repo.FindByMerchantIDForUpdate(ctx, tx, merchantID)
		if err != nil {
			return err
		}
		if wallet.Balance < amount {
			return ErrInsufficientBalance
		}
		wallet.Balance -= amount
		if err := s.repo.Save(ctx, tx, wallet); err != nil {
			return err
		}
		return s.repo.CreateTransaction(ctx, tx, &entity.WalletTransaction{
			WalletID: wallet.ID, SettlementID: settlementID, Type: entity.TxDebit, Amount: amount, Balance: wallet.Balance, Description: "settlement withdrawal",
		})
	})
}
