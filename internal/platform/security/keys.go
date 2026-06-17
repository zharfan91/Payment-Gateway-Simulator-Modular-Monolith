package security

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomToken(prefix string, bytesLen int) (string, error) {
	raw := make([]byte, bytesLen)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(raw), nil
}
