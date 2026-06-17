package service

import "testing"

func TestCalculateFee(t *testing.T) {
	fee, net := CalculateFee(500000)
	if fee != 14500 {
		t.Fatalf("fee = %d, want 14500", fee)
	}
	if net != 485500 {
		t.Fatalf("net = %d, want 485500", net)
	}
}
