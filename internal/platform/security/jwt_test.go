package security

import (
	"testing"
	"time"
)

func TestGenerateAndParseJWT(t *testing.T) {
	token, jti, err := GenerateJWT("user-1", "user@example.com", "access", "secret", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseJWT(token, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user-1" || claims.Email != "user@example.com" || claims.Type != "access" || claims.ID != jti {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
