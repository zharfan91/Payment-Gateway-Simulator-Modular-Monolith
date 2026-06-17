package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName        string
	Env            string
	HTTPAddr       string
	DatabaseDSN    string
	RedisAddr      string
	RedisPassword  string
	RabbitMQURL    string
	AccessSecret   string
	RefreshSecret  string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	PaymentBaseURL string
	JaegerEndpoint string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		AppName:        env("APP_NAME", "payment-gateway-simulator"),
		Env:            env("APP_ENV", "local"),
		HTTPAddr:       env("HTTP_ADDR", ":8080"),
		DatabaseDSN:    env("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=payment_gateway port=5432 sslmode=disable TimeZone=UTC"),
		RedisAddr:      env("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  env("REDIS_PASSWORD", ""),
		RabbitMQURL:    env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		AccessSecret:   env("JWT_ACCESS_SECRET", "local-access-secret"),
		RefreshSecret:  env("JWT_REFRESH_SECRET", "local-refresh-secret"),
		AccessTTL:      time.Duration(envInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
		RefreshTTL:     time.Duration(envInt("JWT_REFRESH_TTL_HOURS", 720)) * time.Hour,
		PaymentBaseURL: env("PAYMENT_BASE_URL", "http://localhost:8080/pay"),
		JaegerEndpoint: env("OTEL_EXPORTER_JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
