package validator

import (
	"errors"
	"strings"

	"github.com/zharf/payment-gateway-simulator/modules/settlement/dto"
)

func Create(req dto.CreateSettlementRequest) error {
	if req.Amount <= 0 || strings.TrimSpace(req.BankAccount) == "" {
		return errors.New("amount and bank_account are required")
	}
	return nil
}
