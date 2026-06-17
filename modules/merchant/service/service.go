package service

import (
	"context"
	"errors"

	"github.com/zharf/payment-gateway-simulator/internal/platform/security"
	"github.com/zharf/payment-gateway-simulator/modules/merchant/dto"
	"github.com/zharf/payment-gateway-simulator/modules/merchant/entity"
	"github.com/zharf/payment-gateway-simulator/modules/merchant/repository"
	"github.com/zharf/payment-gateway-simulator/modules/merchant/validator"
	notificationvalidator "github.com/zharf/payment-gateway-simulator/modules/notification/validator"
	wallet "github.com/zharf/payment-gateway-simulator/modules/wallet/service"
	"gorm.io/gorm"
)

type Service struct {
	repo   *repository.Repository
	wallet *wallet.Service
}

func New(repo *repository.Repository, wallet *wallet.Service) *Service {
	return &Service{repo: repo, wallet: wallet}
}

func (s *Service) Create(ctx context.Context, userID string, req dto.CreateMerchantRequest) (*entity.Merchant, error) {
	if err := validator.Create(req); err != nil {
		return nil, err
	}
	secret, err := security.RandomToken("sk_", 32)
	if err != nil {
		return nil, err
	}
	merchant := &entity.Merchant{UserID: userID, Name: req.Name, Email: req.Email, WebhookURL: req.WebhookURL, SecretKey: secret}
	if err := s.repo.Create(ctx, merchant); err != nil {
		return nil, err
	}
	return merchant, s.wallet.Ensure(ctx, merchant.ID)
}

func (s *Service) Profile(ctx context.Context, userID string) (*entity.Merchant, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) GenerateAPIKey(ctx context.Context, userID string) (string, error) {
	merchant, err := s.Profile(ctx, userID)
	if err != nil {
		return "", err
	}
	key, err := security.RandomToken("pk_", 24)
	if err != nil {
		return "", err
	}
	return key, s.repo.CreateAPIKey(ctx, &entity.APIKey{MerchantID: merchant.ID, Key: key, Active: true})
}

func (s *Service) GenerateSecretKey(ctx context.Context, userID string) (string, error) {
	merchant, err := s.Profile(ctx, userID)
	if err != nil {
		return "", err
	}
	secret, err := security.RandomToken("sk_", 32)
	if err != nil {
		return "", err
	}
	merchant.SecretKey = secret
	return secret, s.repo.Update(ctx, merchant)
}

func (s *Service) RegisterWebhook(ctx context.Context, userID, url string) error {
	if err := notificationvalidator.WebhookURL(url); err != nil {
		return err
	}
	merchant, err := s.Profile(ctx, userID)
	if err != nil {
		return err
	}
	merchant.WebhookURL = url
	return s.repo.Update(ctx, merchant)
}

func (s *Service) MerchantFromAPIKey(ctx context.Context, apiKey string) (*entity.Merchant, error) {
	key, err := s.repo.FindAPIKey(ctx, apiKey)
	if err != nil {
		return nil, errors.New("invalid api key")
	}
	return s.repo.FindByID(ctx, key.MerchantID)
}

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
