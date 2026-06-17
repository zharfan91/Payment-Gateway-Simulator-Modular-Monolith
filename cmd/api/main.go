package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/zharf/payment-gateway-simulator/docs/swagger"
	"github.com/zharf/payment-gateway-simulator/internal/platform/cache"
	"github.com/zharf/payment-gateway-simulator/internal/platform/config"
	"github.com/zharf/payment-gateway-simulator/internal/platform/database"
	"github.com/zharf/payment-gateway-simulator/internal/platform/messaging"
	pgmiddleware "github.com/zharf/payment-gateway-simulator/internal/platform/middleware"
	"github.com/zharf/payment-gateway-simulator/internal/platform/observability"
	audithandler "github.com/zharf/payment-gateway-simulator/modules/auditlog/handler"
	auditrepo "github.com/zharf/payment-gateway-simulator/modules/auditlog/repository"
	auditservice "github.com/zharf/payment-gateway-simulator/modules/auditlog/service"
	authhandler "github.com/zharf/payment-gateway-simulator/modules/auth/handler"
	authrepo "github.com/zharf/payment-gateway-simulator/modules/auth/repository"
	authservice "github.com/zharf/payment-gateway-simulator/modules/auth/service"
	merchanthandler "github.com/zharf/payment-gateway-simulator/modules/merchant/handler"
	merchantrepo "github.com/zharf/payment-gateway-simulator/modules/merchant/repository"
	merchantservice "github.com/zharf/payment-gateway-simulator/modules/merchant/service"
	notificationhandler "github.com/zharf/payment-gateway-simulator/modules/notification/handler"
	notificationrepo "github.com/zharf/payment-gateway-simulator/modules/notification/repository"
	notificationservice "github.com/zharf/payment-gateway-simulator/modules/notification/service"
	paymenthandler "github.com/zharf/payment-gateway-simulator/modules/payment/handler"
	paymentrepo "github.com/zharf/payment-gateway-simulator/modules/payment/repository"
	paymentservice "github.com/zharf/payment-gateway-simulator/modules/payment/service"
	settlementhandler "github.com/zharf/payment-gateway-simulator/modules/settlement/handler"
	settlementrepo "github.com/zharf/payment-gateway-simulator/modules/settlement/repository"
	settlementservice "github.com/zharf/payment-gateway-simulator/modules/settlement/service"
	wallethandler "github.com/zharf/payment-gateway-simulator/modules/wallet/handler"
	walletrepo "github.com/zharf/payment-gateway-simulator/modules/wallet/repository"
	walletservice "github.com/zharf/payment-gateway-simulator/modules/wallet/service"
)

// @title Payment Gateway Simulator API
// @version 1.0
// @description Modular monolith payment gateway simulator.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	ctx := context.Background()
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	redis, err := cache.Connect(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		slog.Warn("redis unavailable; cache-backed features degraded", "error", err)
		redis = nil
	}
	var publisher messaging.Publisher = messaging.NoopPublisher{}
	if mq, err := messaging.Connect(cfg.RabbitMQURL); err == nil {
		defer mq.Close()
		publisher = mq
	} else {
		slog.Warn("rabbitmq unavailable; event publishing disabled", "error", err)
	}
	if shutdown, err := observability.InitTracing(ctx, cfg.AppName, cfg.JaegerEndpoint); err == nil {
		defer shutdown(ctx)
	}

	auditRepo := auditrepo.New(db)
	auditSvc := auditservice.New(auditRepo)
	walletSvc := walletservice.New(walletrepo.New(db))
	merchantRepo := merchantrepo.New(db)
	merchantSvc := merchantservice.New(merchantRepo, walletSvc)
	authSvc := authservice.New(cfg, authrepo.New(db), redis, auditSvc)
	paymentSvc := paymentservice.New(cfg, paymentrepo.New(db), redis, publisher, walletSvc)
	settlementSvc := settlementservice.New(settlementrepo.New(db), walletSvc)
	notificationSvc := notificationservice.New(notificationrepo.New(db), merchantRepo)

	app := fiber.New(fiber.Config{AppName: cfg.AppName, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second})
	app.Use(recover.New(), pgmiddleware.RequestID(), pgmiddleware.RateLimit(redis, 120, time.Minute))
	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	app.Get("/ready", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ready"}) })
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api/v1")
	jwt := pgmiddleware.JWT(cfg, redis)
	authhandler.New(authSvc).RegisterRoutes(api, jwt)
	merchanthandler.New(merchantSvc).RegisterRoutes(api, jwt)
	paymenthandler.New(paymentSvc, merchantSvc).RegisterRoutes(api, jwt)
	wallethandler.New(walletSvc, merchantSvc).RegisterRoutes(api, jwt)
	settlementhandler.New(settlementSvc, merchantSvc).RegisterRoutes(api, jwt)
	notificationhandler.New(notificationSvc, merchantSvc).RegisterRoutes(api, jwt)
	audithandler.New(auditRepo, merchantSvc).RegisterRoutes(api, jwt)

	log.Fatal(app.Listen(cfg.HTTPAddr))
}
