package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/zharf/payment-gateway-simulator/internal/platform/security"
	merchantrepo "github.com/zharf/payment-gateway-simulator/modules/merchant/repository"
	"github.com/zharf/payment-gateway-simulator/modules/notification/entity"
	"github.com/zharf/payment-gateway-simulator/modules/notification/repository"
)

type Service struct {
	repo         *repository.Repository
	merchantRepo *merchantrepo.Repository
	client       *http.Client
}

func New(repo *repository.Repository, merchantRepo *merchantrepo.Repository) *Service {
	return &Service{repo: repo, merchantRepo: merchantRepo, client: &http.Client{Timeout: 5 * time.Second}}
}

func (s *Service) SendWebhook(ctx context.Context, merchantID, event string, payload any) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil || merchant.WebhookURL == "" {
		return
	}
	body, _ := json.Marshal(map[string]any{"event": event, "data": payload})
	for attempt := 1; attempt <= 3; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, merchant.WebhookURL, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Signature", security.SignHMAC(merchant.SecretKey, body))
		resp, err := s.client.Do(req)
		status := 0
		success := false
		errText := ""
		if resp != nil {
			status = resp.StatusCode
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			success = status >= 200 && status < 300
		}
		if err != nil {
			errText = err.Error()
		}
		_ = s.repo.CreateWebhookLog(ctx, &entity.WebhookLog{MerchantID: merchantID, Event: event, URL: merchant.WebhookURL, Payload: string(body), StatusCode: status, Attempt: attempt, Success: success, Error: errText})
		if success {
			return
		}
		time.Sleep(time.Duration(attempt) * 250 * time.Millisecond)
	}
}
