package service

import (
	"context"
	"encoding/json"

	"github.com/zharf/payment-gateway-simulator/modules/auditlog/entity"
	"github.com/zharf/payment-gateway-simulator/modules/auditlog/repository"
)

type Service struct{ repo *repository.Repository }

func New(repo *repository.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Record(ctx context.Context, actorID, merchantID, action, resource string, metadata any) {
	body, _ := json.Marshal(metadata)
	_ = s.repo.Create(ctx, &entity.AuditLog{
		ActorID: actorID, MerchantID: merchantID, Action: action, Resource: resource, Metadata: string(body),
	})
}
