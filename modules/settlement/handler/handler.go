package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	merchant "github.com/zharf/payment-gateway-simulator/modules/merchant/service"
	"github.com/zharf/payment-gateway-simulator/modules/settlement/dto"
	"github.com/zharf/payment-gateway-simulator/modules/settlement/service"
)

type Handler struct {
	svc      *service.Service
	merchant *merchant.Service
}

func New(svc *service.Service, merchant *merchant.Service) *Handler {
	return &Handler{svc: svc, merchant: merchant}
}

func (h *Handler) RegisterRoutes(r fiber.Router, auth fiber.Handler) {
	g := r.Group("/settlements", auth)
	g.Post("/", h.Create)
	g.Get("/", h.List)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req dto.CreateSettlementRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	merchantEntity, err := h.merchant.Profile(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "merchant not found")
	}
	settlement, err := h.svc.Create(c.Context(), merchantEntity.ID, req)
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.Created(c, settlement)
}

func (h *Handler) List(c *fiber.Ctx) error {
	merchantEntity, err := h.merchant.Profile(c.Context(), c.Locals("user_id").(string))
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "merchant not found")
	}
	settlements, err := h.svc.List(c.Context(), merchantEntity.ID)
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.OK(c, settlements)
}
