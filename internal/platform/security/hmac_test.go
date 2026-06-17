package security

import "testing"

func TestVerifyHMAC(t *testing.T) {
	body := []byte(`{"amount":500000}`)
	signature := SignHMAC("secret", body)
	if !VerifyHMAC("secret", body, signature) {
		t.Fatal("expected valid signature")
	}
	if VerifyHMAC("secret", []byte(`{"amount":1}`), signature) {
		t.Fatal("expected invalid signature for modified body")
	}
}
