package pulse

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

// ComputeWebhookSignature calculates the HMAC-SHA256 hex string digest for a webhook payload.
func ComputeWebhookSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyWebhookSignature cryptographically validates a webhook signature header against the payload and secret.
// Utilizes constant-time comparison to guard against timing side-channel attacks.
func VerifyWebhookSignature(payload []byte, signature, secret string) bool {
	if strings.TrimSpace(signature) == "" || strings.TrimSpace(secret) == "" {
		return false
	}

	expectedHex := ComputeWebhookSignature(payload, secret)
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	expectedBytes, err := hex.DecodeString(expectedHex)
	if err != nil {
		return false
	}

	return hmac.Equal(sigBytes, expectedBytes)
}

// WebhookMiddleware provides standard net/http middleware to protect webhook endpoints.
// Verifies the "X-Pulse-Signature" header against the raw body bytes and configured secret.
func WebhookMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sig := r.Header.Get("X-Pulse-Signature")
			if sig == "" {
				http.Error(w, `{"error":"missing X-Pulse-Signature header"}`, http.StatusUnauthorized)
				return
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, `{"error":"unable to read request body"}`, http.StatusBadRequest)
				return
			}
			// Restore the body for downstream handlers
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			if !VerifyWebhookSignature(bodyBytes, sig, secret) {
				http.Error(w, `{"error":"invalid webhook signature"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

