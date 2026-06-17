package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	"github.com/zharf/payment-gateway-simulator/modules/auth/dto"
	"github.com/zharf/payment-gateway-simulator/modules/auth/service"
)

type Handler struct{ svc *service.Service }

func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(r fiber.Router, auth fiber.Handler) {
	g := r.Group("/auth")
	g.Post("/register", h.Register)
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
	g.Post("/logout", auth, h.Logout)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	res, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.Created(c, res)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	res, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return httpx.Error(c, fiber.StatusUnauthorized, err.Error())
	}
	return httpx.OK(c, res)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	res, err := h.svc.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return httpx.Error(c, fiber.StatusUnauthorized, err.Error())
	}
	return httpx.OK(c, res)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	var req dto.LogoutRequest
	_ = c.BodyParser(&req)
	if err := h.svc.Logout(c.Context(), c.Locals("jti").(string), c.Locals("user_id").(string), req.RefreshToken); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return httpx.OK(c, fiber.Map{"logged_out": true})
}
