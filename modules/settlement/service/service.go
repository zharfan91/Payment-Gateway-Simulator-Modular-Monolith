package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zharf/payment-gateway-simulator/modules/settlement/dto"
	"github.com/zharf/payment-gateway-simulator/modules/settlement/entity"
	"github.com/zharf/payment-gateway-simulator/modules/settlement/repository"
	"github.com/zharf/payment-gateway-simulator/modules/settlement/validator"
	wallet "github.com/zharf/payment-gateway-simulator/modules/wallet/service"
)

type Service struct {
	repo   *repository.Repository
	wallet *wallet.Service
}

func New(repo *repository.Repository, wallet *wallet.Service) *Service {
	return &Service{repo: repo, wallet: wallet}
}

func (s *Service) Create(ctx context.Context, merchantID string, req dto.CreateSettlementRequest) (*entity.Settlement, error) {
	if err := validator.Create(req); err != nil {
		return nil, err
	}
	currentWallet, err := s.wallet.Get(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	if currentWallet.Balance < req.Amount {
		return nil, wallet.ErrInsufficientBalance
	}
	settlement := &entity.Settlement{ID: uuid.NewString(), MerchantID: merchantID, Amount: req.Amount, BankAccount: req.BankAccount, Status: entity.StatusPaid}
	if err := s.repo.Create(ctx, settlement); err != nil {
		return nil, err
	}
	if err := s.wallet.Debit(ctx, merchantID, settlement.ID, req.Amount); err != nil {
		return nil, err
	}
	return settlement, nil
}

func (s *Service) List(ctx context.Context, merchantID string) ([]entity.Settlement, error) {
	return s.repo.ListByMerchant(ctx, merchantID)
}
