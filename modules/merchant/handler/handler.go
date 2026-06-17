package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	"github.com/zharf/payment-gateway-simulator/modules/merchant/dto"
	"github.com/zharf/payment-gateway-simulator/modules/merchant/service"
)

type Handler struct{ svc *service.Service }

func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(r fiber.Router, auth fiber.Handler) {
	g := r.Group("/merchants", auth)
	g.Post("/", h.Create)
	g.Get("/me", h.Profile)
	g.Post("/api-keys", h.GenerateAPIKey)
	g.Post("/secret-keys", h.GenerateSecretKey)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req dto.CreateMerchantRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	merchant, err := h.svc.Create(c.Context(), c.Locals("user_id").(string), req)
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.Created(c, dto.MerchantResponse{ID: merchant.ID, Name: merchant.Name, Email: merchant.Email, WebhookURL: merchant.WebhookURL})
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	merchant, err := h.svc.Profile(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "merchant not found")
	}
	return httpx.OK(c, dto.MerchantResponse{ID: merchant.ID, Name: merchant.Name, Email: merchant.Email, WebhookURL: merchant.WebhookURL})
}

func (h *Handler) GenerateAPIKey(c *fiber.Ctx) error {
	key, err := h.svc.GenerateAPIKey(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.Created(c, dto.APIKeyResponse{APIKey: key})
}

func (h *Handler) GenerateSecretKey(c *fiber.Ctx) error {
	secret, err := h.svc.GenerateSecretKey(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.Created(c, dto.SecretKeyResponse{SecretKey: secret})
}
