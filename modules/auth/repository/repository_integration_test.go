//go:build integration

package repository_test

import (
	"context"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zharf/payment-gateway-simulator/modules/auth/entity"
	"github.com/zharf/payment-gateway-simulator/modules/auth/repository"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepositoryPostgres(t *testing.T) {
	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("payment_gateway_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	if err != nil {
		if strings.Contains(err.Error(), "Docker") || strings.Contains(err.Error(), "docker") {
			t.Skipf("docker provider unavailable: %v", err)
		}
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		t.Fatal(err)
	}
	repo := repository.New(db)
	user := &entity.User{Name: "Jane", Email: "jane@example.com", PasswordHash: "hash"}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatal(err)
	}
	found, err := repo.FindByEmail(ctx, "jane@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if found.ID == "" || found.Email != user.Email {
		t.Fatalf("unexpected user: %+v", found)
	}
}
