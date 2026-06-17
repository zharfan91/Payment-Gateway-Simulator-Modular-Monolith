package validator

import (
	"errors"
	"strings"

	"github.com/zharf/payment-gateway-simulator/modules/auth/dto"
)

func Register(req dto.RegisterRequest) error {
	if strings.TrimSpace(req.Name) == "" || !strings.Contains(req.Email, "@") || len(req.Password) < 8 {
		return errors.New("name, valid email, and password with at least 8 chars are required")
	}
	return nil
}

func Login(req dto.LoginRequest) error {
	if !strings.Contains(req.Email, "@") || req.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}
