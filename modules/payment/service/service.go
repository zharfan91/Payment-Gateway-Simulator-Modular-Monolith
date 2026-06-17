package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/zharf/payment-gateway-simulator/internal/platform/config"
	"github.com/zharf/payment-gateway-simulator/internal/platform/messaging"
	"github.com/zharf/payment-gateway-simulator/internal/platform/security"
	"github.com/zharf/payment-gateway-simulator/modules/payment/dto"
	"github.com/zharf/payment-gateway-simulator/modules/payment/entity"
	"github.com/zharf/payment-gateway-simulator/modules/payment/repository"
	"github.com/zharf/payment-gateway-simulator/modules/payment/validator"
	wallet "github.com/zharf/payment-gateway-simulator/modules/wallet/service"
	"gorm.io/gorm"
)

type Service struct {
	cfg       config.Config
	repo      *repository.Repository
	redis     *redis.Client
	publisher messaging.Publisher
	wallet    *wallet.Service
}

func New(cfg config.Config, repo *repository.Repository, redis *redis.Client, publisher messaging.Publisher, wallet *wallet.Service) *Service {
	return &Service{cfg: cfg, repo: repo, redis: redis, publisher: publisher, wallet: wallet}
}

func (s *Service) Create(ctx context.Context, merchantID, idempotencyKey string, req dto.CreatePaymentRequest) (*entity.Payment, error) {
	if err := validator.Create(req); err != nil {
		return nil, err
	}
	if idempotencyKey != "" {
		if cached, err := s.fromIdempotency(ctx, merchantID, idempotencyKey); err == nil {
			return cached, nil
		}
	}
	fee, net := wallet.CalculateFee(req.Amount)
	token, err := security.RandomToken("pay_", 24)
	if err != nil {
		return nil, err
	}
	payment := &entity.Payment{
		MerchantID: merchantID, ExternalID: req.ExternalID, Amount: req.Amount, Fee: fee, NetAmount: net,
		Status: entity.StatusPending, PaymentToken: token, PaymentURL: s.cfg.PaymentBaseURL + "/" + token,
		Description: req.Description, IdempotencyKey: idempotencyKey,
	}
	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}
	_ = s.repo.CreateLog(ctx, &entity.PaymentLog{PaymentID: payment.ID, Status: payment.Status, Message: "payment created"})
	_ = s.publisher.Publish(ctx, "payment.created", payment)
	s.cachePayment(ctx, payment)
	return payment, nil
}

func (s *Service) Get(ctx context.Context, id string) (*entity.Payment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) MarkSuccess(ctx context.Context, id string) (*entity.Payment, error) {
	payment, err := s.transition(ctx, id, entity.StatusSuccess)
	if err != nil {
		return nil, err
	}
	if err := s.wallet.Credit(ctx, payment.MerchantID, payment.ID, payment.NetAmount); err != nil {
		return nil, err
	}
	_ = s.publisher.Publish(ctx, "payment.success", payment)
	return payment, nil
}

func (s *Service) MarkFailed(ctx context.Context, id string) (*entity.Payment, error) {
	payment, err := s.transition(ctx, id, entity.StatusFailed)
	if err == nil {
		_ = s.publisher.Publish(ctx, "payment.failed", payment)
	}
	return payment, err
}

func (s *Service) Refund(ctx context.Context, id string) (*entity.Payment, error) {
	payment, err := s.transition(ctx, id, entity.StatusRefunded)
	if err == nil {
		_ = s.publisher.Publish(ctx, "payment.refunded", payment)
	}
	return payment, err
}

func (s *Service) transition(ctx context.Context, id string, next entity.Status) (*entity.Payment, error) {
	payment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment.Status != entity.StatusPending && !(payment.Status == entity.StatusSuccess && next == entity.StatusRefunded) {
		return nil, errors.New("invalid payment status transition")
	}
	payment.Status = next
	if err := s.repo.Update(ctx, payment); err != nil {
		return nil, err
	}
	_ = s.repo.CreateLog(ctx, &entity.PaymentLog{PaymentID: payment.ID, Status: next, Message: "status updated"})
	s.cachePayment(ctx, payment)
	return payment, nil
}

func (s *Service) fromIdempotency(ctx context.Context, merchantID, key string) (*entity.Payment, error) {
	if s.redis != nil {
		var payment entity.Payment
		value, err := s.redis.Get(ctx, "idempotency:"+merchantID+":"+key).Bytes()
		if err == nil && json.Unmarshal(value, &payment) == nil {
			return &payment, nil
		}
	}
	return s.repo.FindByIdempotencyKey(ctx, merchantID, key)
}

func (s *Service) cachePayment(ctx context.Context, payment *entity.Payment) {
	if s.redis == nil {
		return
	}
	body, _ := json.Marshal(payment)
	_ = s.redis.Set(ctx, "payment:"+payment.ID, body, 0).Err()
	if payment.IdempotencyKey != "" {
		_ = s.redis.Set(ctx, "idempotency:"+payment.MerchantID+":"+payment.IdempotencyKey, body, 0).Err()
	}
}

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
