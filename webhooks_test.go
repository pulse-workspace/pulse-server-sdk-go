package pulse

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhooks(t *testing.T) {
	secret := "webhook_secret_key_123456789"
	payload := []byte(`{"event":"conversation.created","convId":"conv_123"}`)

	t.Run("compute and verify valid signature", func(t *testing.T) {
		sig := ComputeWebhookSignature(payload, secret)
		if len(sig) != 64 {
			t.Errorf("expected 64 hex characters, got %d", len(sig))
		}

		if !VerifyWebhookSignature(payload, sig, secret) {
			t.Errorf("expected signature to be verified")
		}
	})

	t.Run("altered payload fails verification", func(t *testing.T) {
		sig := ComputeWebhookSignature(payload, secret)
		altered := []byte(`{"event":"conversation.created","convId":"conv_999"}`)

		if VerifyWebhookSignature(altered, sig, secret) {
			t.Errorf("expected altered payload to fail verification")
		}
	})

	t.Run("middleware verification", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		middleware := WebhookMiddleware(secret)(handler)

		// Case 1: missing header -> 401
		req1 := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payload))
		rec1 := httptest.NewRecorder()
		middleware.ServeHTTP(rec1, req1)
		if rec1.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 for missing signature, got %d", rec1.Code)
		}

		// Case 2: valid signature -> 200
		sig := ComputeWebhookSignature(payload, secret)
		req2 := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payload))
		req2.Header.Set("X-Pulse-Signature", sig)
		rec2 := httptest.NewRecorder()
		middleware.ServeHTTP(rec2, req2)
		if rec2.Code != http.StatusOK {
			t.Errorf("expected 200 for valid signature, got %d", rec2.Code)
		}
	})
}

