package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zharf/payment-gateway-simulator/internal/platform/config"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
	"github.com/zharf/payment-gateway-simulator/internal/platform/security"
)

func JWT(cfg config.Config, redis *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return httpx.Error(c, fiber.StatusUnauthorized, "missing bearer token")
		}
		claims, err := security.ParseJWT(strings.TrimPrefix(header, "Bearer "), cfg.AccessSecret)
		if err != nil || claims.Type != "access" {
			return httpx.Error(c, fiber.StatusUnauthorized, "invalid token")
		}
		if redis != nil {
			blacklisted, _ := redis.Exists(context.Background(), "jwt:blacklist:"+claims.ID).Result()
			if blacklisted > 0 {
				return httpx.Error(c, fiber.StatusUnauthorized, "token revoked")
			}
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("jti", claims.ID)
		return c.Next()
	}
}
