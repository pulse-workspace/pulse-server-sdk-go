# Pulse Server-Side SDK for Go (`pulse-server-sdk-go`)

Official server-side Golang SDK for the **Pulse Real-time Messaging Platform**.

Enables Go backend services, microservices, and worker daemons to mint user JWT session tokens, manage channels and users, dispatch system messages, and verify inbound webhook signatures.

## 📦 Installation

```bash
go get github.com/pulse-platform/pulse-server-sdk-go
```

## 🚀 Quick Start

### 1. Initialize Client

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/pulse-platform/pulse-server-sdk-go"
)

func main() {
	client, err := pulse.NewClient(pulse.Config{
		BaseURL:    "http://localhost:32011",
		AppID:      "app_production",
		AppSecret:  "sec_0123456789abcdef",
		Timeout:    5 * time.Second,
		MaxRetries: 3,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Pulse client: %v", err)
	}

	// 2. Mint client JWT session token for WebSocket
	token, err := client.Auth.CreateToken("usr_alice", pulse.TokenOptions{
		Role: "user",
		TTL:  24 * time.Hour,
	})
	if err != nil {
		log.Fatalf("Failed to create token: %v", err)
	}
	log.Printf("Generated client token: %s", token)

	// 3. Create group conversation
	ctx := context.Background()
	conv, err := client.Conversations.Create(ctx, pulse.CreateConversationRequest{
		Type:      "GROUP",
		Name:      "Golang Backend Core",
		MemberIDs: []string{"usr_alice", "usr_bob"},
	})
	if err != nil {
		log.Fatalf("Failed to create conversation: %v", err)
	}

	// 4. Send official system message
	_, err = client.Messages.SendSystemMessage(ctx, pulse.SystemMessageRequest{
		ConvID: conv.ConvID,
		Body:   "Channel created successfully via Go Server SDK.",
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
```

### 5. Webhook Middleware Integration

```go
package main

import (
	"net/http"

	"github.com/pulse-platform/pulse-server-sdk-go"
)

func main() {
	webhookHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"received":true}`))
	})

	// Wrap handler with automatic HMAC-SHA256 signature verification
	http.Handle("/webhooks/pulse", pulse.WebhookMiddleware("sec_0123456789abcdef")(webhookHandler))
	http.ListenAndServe(":8080", nil)
}
```

