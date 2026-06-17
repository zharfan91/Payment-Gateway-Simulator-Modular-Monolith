package validator

import (
	"errors"

	"github.com/zharf/payment-gateway-simulator/modules/payment/dto"
)

func Create(req dto.CreatePaymentRequest) error {
	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	return nil
}
