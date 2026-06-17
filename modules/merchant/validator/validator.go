package validator

import (
	"errors"
	"strings"

	"github.com/zharf/payment-gateway-simulator/modules/merchant/dto"
)

func Create(req dto.CreateMerchantRequest) error {
	if strings.TrimSpace(req.Name) == "" || !strings.Contains(req.Email, "@") {
		return errors.New("merchant name and valid email are required")
	}
	return nil
}
