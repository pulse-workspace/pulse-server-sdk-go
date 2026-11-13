package pulse

import (
	"strings"
	"testing"
	"time"
)

func TestAuthService_CreateAndVerifyToken(t *testing.T) {
	appID := "app_test_tenant"
	appSecret := "0123456789abcdef0123456789abcdef"
	auth := NewAuthService(appID, appSecret)

	t.Run("create valid token", func(t *testing.T) {
		token, err := auth.CreateToken("usr_alice", TokenOptions{
			Role: "admin",
			TTL:  2 * time.Hour,
			Metadata: map[string]interface{}{
				"tier": "enterprise",
			},
		})
		if err != nil {
			t.Fatalf("unexpected error creating token: %v", err)
		}

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			t.Fatalf("expected 3 parts in JWT, got %d", len(parts))
		}

		claims, err := auth.VerifyToken(token)
		if err != nil {
			t.Fatalf("unexpected error verifying token: %v", err)
		}

		if claims["sub"] != "usr_alice" {
			t.Errorf("expected sub usr_alice, got %v", claims["sub"])
		}
		if claims["iss"] != appID {
			t.Errorf("expected iss %s, got %v", appID, claims["iss"])
		}
		if claims["role"] != "admin" {
			t.Errorf("expected role admin, got %v", claims["role"])
		}
	})

	t.Run("empty user ID fails", func(t *testing.T) {
		_, err := auth.CreateToken("")
		if err != ErrEmptyUserID {
			t.Fatalf("expected ErrEmptyUserID, got %v", err)
		}
	})

	t.Run("tampered signature fails", func(t *testing.T) {
		token, err := auth.CreateToken("usr_bob")
		if err != nil {
			t.Fatal(err)
		}

		parts := strings.Split(token, ".")
		tampered := parts[0] + "." + parts[1] + ".invalid_sig"

		_, err = auth.VerifyToken(tampered)
		if err != ErrInvalidSignature {
			t.Fatalf("expected ErrInvalidSignature, got %v", err)
		}
	})

	t.Run("foreign secret fails", func(t *testing.T) {
		foreignAuth := NewAuthService(appID, "another_secret_key_123456789")
		token, err := foreignAuth.CreateToken("usr_charlie")
		if err != nil {
			t.Fatal(err)
		}

		_, err = auth.VerifyToken(token)
		if err != ErrInvalidSignature {
			t.Fatalf("expected ErrInvalidSignature, got %v", err)
		}
	})
}

