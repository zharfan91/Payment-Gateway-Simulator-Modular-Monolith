package httpx

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Error     string `json:"error"`
}

func OK(c *fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{"data": data})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": data})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{RequestID: c.GetRespHeader("X-Request-ID"), Error: message})
}
