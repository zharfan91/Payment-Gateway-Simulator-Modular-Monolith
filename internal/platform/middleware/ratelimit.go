package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zharf/payment-gateway-simulator/internal/platform/httpx"
)

func RateLimit(redis *redis.Client, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if redis == nil {
			return c.Next()
		}
		key := "rate:" + c.IP()
		ctx := context.Background()
		count, err := redis.Incr(ctx, key).Result()
		if err == nil && count == 1 {
			_ = redis.Expire(ctx, key, window).Err()
		}
		if err == nil && count > int64(limit) {
			return httpx.Error(c, fiber.StatusTooManyRequests, "rate limit exceeded")
		}
		return c.Next()
	}
}
