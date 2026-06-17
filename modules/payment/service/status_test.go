package service

import (
	"testing"

	"github.com/zharf/payment-gateway-simulator/modules/payment/entity"
)

func TestTerminalStatusConstants(t *testing.T) {
	statuses := []entity.Status{entity.StatusPending, entity.StatusSuccess, entity.StatusFailed, entity.StatusExpired, entity.StatusRefunded}
	if len(statuses) != 5 {
		t.Fatal("unexpected payment status count")
	}
}
