package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	merchant "github.com/zharf/payment-gateway-simulator/modules/merchant/service"
	"github.com/zharf/payment-gateway-simulator/modules/notification/dto"
	"github.com/zharf/payment-gateway-simulator/modules/notification/service"
	"github.com/zharf/payment-gateway-simulator/modules/notification/validator"
)

type Handler struct {
	svc      *service.Service
	merchant *merchant.Service
}

func New(svc *service.Service, merchant *merchant.Service) *Handler {
	return &Handler{svc: svc, merchant: merchant}
}

func (h *Handler) RegisterRoutes(r fiber.Router, auth fiber.Handler) {
	g := r.Group("/webhooks", auth)
	g.Post("/", h.RegisterWebhook)
	g.Post("/test", h.TestWebhook)
}

func (h *Handler) RegisterWebhook(c *fiber.Ctx) error {
	var req dto.RegisterWebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	if err := validator.WebhookURL(req.WebhookURL); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if err := h.merchant.RegisterWebhook(c.Context(), c.Locals("user_id").(string), req.WebhookURL); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.OK(c, fiber.Map{"registered": true})
}

func (h *Handler) TestWebhook(c *fiber.Ctx) error {
	var req dto.TestWebhookRequest
	_ = c.BodyParser(&req)
	if req.Event == "" {
		req.Event = "webhook.test"
	}
	merchantEntity, err := h.merchant.Profile(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "merchant not found")
	}
	go h.svc.SendWebhook(c.Context(), merchantEntity.ID, req.Event, fiber.Map{"ok": true})
	return httpx.OK(c, fiber.Map{"queued": true})
}
