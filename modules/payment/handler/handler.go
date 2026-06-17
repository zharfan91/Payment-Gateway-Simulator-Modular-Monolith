package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	"github.com/zharf/payment-gateway-simulator/internal/platform/security"
	merchant "github.com/zharf/payment-gateway-simulator/modules/merchant/service"
	"github.com/zharf/payment-gateway-simulator/modules/payment/dto"
	"github.com/zharf/payment-gateway-simulator/modules/payment/entity"
	"github.com/zharf/payment-gateway-simulator/modules/payment/service"
)

type Handler struct {
	svc      *service.Service
	merchant *merchant.Service
}

func New(svc *service.Service, merchant *merchant.Service) *Handler {
	return &Handler{svc: svc, merchant: merchant}
}

func (h *Handler) RegisterRoutes(r fiber.Router, auth fiber.Handler) {
	g := r.Group("/payments")
	g.Post("/", h.Create)
	g.Get("/:id", auth, h.Get)
	g.Post("/:id/simulate-success", auth, h.Success)
	g.Post("/:id/simulate-failed", auth, h.Failed)
	g.Post("/:id/refund", auth, h.Refund)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	merchantEntity, err := h.merchant.MerchantFromAPIKey(c.Context(), c.Get("X-API-Key"))
	if err != nil {
		return httpx.Error(c, fiber.StatusUnauthorized, "invalid api key")
	}
	if sig := c.Get("X-Signature"); sig != "" && !security.VerifyHMAC(merchantEntity.SecretKey, c.BodyRaw(), sig) {
		return httpx.Error(c, fiber.StatusUnauthorized, "invalid signature")
	}
	var req dto.CreatePaymentRequest
	if err := json.Unmarshal(c.BodyRaw(), &req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	payment, err := h.svc.Create(c.Context(), merchantEntity.ID, c.Get("Idempotency-Key"), req)
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.Created(c, toResponse(payment))
}

func (h *Handler) Get(c *fiber.Ctx) error {
	payment, err := h.svc.Get(c.Context(), c.Params("id"))
	if err != nil {
		return httpx.Error(c, fiber.StatusNotFound, "payment not found")
	}
	return httpx.OK(c, toResponse(payment))
}

func (h *Handler) Success(c *fiber.Ctx) error {
	payment, err := h.svc.MarkSuccess(c.Context(), c.Params("id"))
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.OK(c, toResponse(payment))
}

func (h *Handler) Failed(c *fiber.Ctx) error {
	payment, err := h.svc.MarkFailed(c.Context(), c.Params("id"))
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.OK(c, toResponse(payment))
}

func (h *Handler) Refund(c *fiber.Ctx) error {
	payment, err := h.svc.Refund(c.Context(), c.Params("id"))
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.OK(c, toResponse(payment))
}

func toResponse(payment *entity.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID: payment.ID, MerchantID: payment.MerchantID, ExternalID: payment.ExternalID, Amount: payment.Amount,
		Fee: payment.Fee, NetAmount: payment.NetAmount, Status: string(payment.Status), PaymentToken: payment.PaymentToken, PaymentURL: payment.PaymentURL,
	}
}
