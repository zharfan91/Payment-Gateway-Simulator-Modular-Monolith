package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	"github.com/zharf/payment-gateway-simulator/modules/auditlog/repository"
	merchant "github.com/zharf/payment-gateway-simulator/modules/merchant/service"
)

type Handler struct {
	repo     *repository.Repository
	merchant *merchant.Service
}

func New(repo *repository.Repository, merchant *merchant.Service) *Handler {
	return &Handler{repo: repo, merchant: merchant}
}

func (h *Handler) RegisterRoutes(r fiber.Router, auth fiber.Handler) {
	r.Get("/audit-logs", auth, h.List)
}

func (h *Handler) List(c *fiber.Ctx) error {
	merchantEntity, err := h.merchant.Profile(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "merchant not found")
	}
	logs, err := h.repo.ListByMerchant(c.Context(), merchantEntity.ID)
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.OK(c, logs)
}
