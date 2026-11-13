package pulse

import (
	"testing"
	"time"
)

func TestClient_Initialization(t *testing.T) {
	cfg := Config{
		BaseURL:   "http://localhost:32011",
		AppID:     "app_dev_tenant",
		AppSecret: "0123456789abcdef0123456789abcdef",
		Timeout:   5 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	if client.Auth == nil {
		t.Error("expected Auth service to be initialized")
	}
	if client.Users == nil {
		t.Error("expected Users service to be initialized")
	}
	if client.Conversations == nil {
		t.Error("expected Conversations service to be initialized")
	}
	if client.Messages == nil {
		t.Error("expected Messages service to be initialized")
	}

	// Verify convenience methods
	raw := []byte("hello world")
	sig := client.ComputeWebhookSignature(raw)
	if !client.VerifyWebhook(raw, sig) {
		t.Error("expected webhook signature to verify")
	}
}

func BenchmarkAuthService_CreateToken(b *testing.B) {
	auth := NewAuthService("app_benchmark", "0123456789abcdef0123456789abcdef")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = auth.CreateToken("usr_bench", TokenOptions{Role: "user"})
	}
}

func BenchmarkWebhook_Verification(b *testing.B) {
	secret := "0123456789abcdef0123456789abcdef"
	payload := []byte(`{"event":"test","timestamp":1700000000}`)
	sig := ComputeWebhookSignature(payload, secret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyWebhookSignature(payload, sig, secret)
	}
}

