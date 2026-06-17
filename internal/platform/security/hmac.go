package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignHMAC(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyHMAC(secret string, body []byte, signature string) bool {
	expected := SignHMAC(secret, body)
	return hmac.Equal([]byte(expected), []byte(signature))
}
