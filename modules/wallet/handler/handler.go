package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	merchant "github.com/zharf/payment-gateway-simulator/modules/merchant/service"
	"github.com/zharf/payment-gateway-simulator/modules/wallet/dto"
	"github.com/zharf/payment-gateway-simulator/modules/wallet/service"
)

type Handler struct {
	svc      *service.Service
	merchant *merchant.Service
}

func New(svc *service.Service, merchant *merchant.Service) *Handler {
	return &Handler{svc: svc, merchant: merchant}
}

func (h *Handler) RegisterRoutes(r fiber.Router, auth fiber.Handler) {
	r.Get("/wallets/me", auth, h.Get)
}

func (h *Handler) Get(c *fiber.Ctx) error {
	merchantEntity, err := h.merchant.Profile(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "merchant not found")
	}
	wallet, err := h.svc.Get(c.Context(), merchantEntity.ID)
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "wallet not found")
	}
	return httpx.OK(c, dto.WalletResponse{MerchantID: merchantEntity.ID, Balance: wallet.Balance})
}
