package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/zharf/payment-gateway-simulator/internal/platform/config"
	"github.com/zharf/payment-gateway-simulator/internal/platform/database"
	"github.com/zharf/payment-gateway-simulator/internal/platform/messaging"
	merchantrepo "github.com/zharf/payment-gateway-simulator/modules/merchant/repository"
	notificationrepo "github.com/zharf/payment-gateway-simulator/modules/notification/repository"
	notificationservice "github.com/zharf/payment-gateway-simulator/modules/notification/service"
	paymententity "github.com/zharf/payment-gateway-simulator/modules/payment/entity"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	mq, err := messaging.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()
	notifications := notificationservice.New(notificationrepo.New(db), merchantrepo.New(db))
	keys := []string{"payment.created", "payment.success", "payment.failed", "payment.refunded"}
	log.Fatal(mq.Consume(context.Background(), "payment-notifications", keys, func(ctx context.Context, event messaging.Event) error {
		var payment paymententity.Payment
		if err := json.Unmarshal(event.Payload, &payment); err != nil {
			return err
		}
		slog.Info("email simulation", "event", event.Name, "payment_id", payment.ID, "merchant_id", payment.MerchantID)
		notifications.SendWebhook(ctx, payment.MerchantID, event.Name, payment)
		return nil
	}))
}
