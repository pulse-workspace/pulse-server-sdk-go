package pulse

import (
	"errors"
	"time"
)

var (
	// ErrNilConfig indicates the configuration struct pointer is nil.
	ErrNilConfig = errors.New("pulse: configuration cannot be nil")
	// ErrMissingBaseURL indicates the BaseURL parameter was omitted.
	ErrMissingBaseURL = errors.New("pulse: BaseURL must be a non-empty string starting with http:// or https://")
	// ErrMissingAppID indicates the AppID parameter was omitted.
	ErrMissingAppID = errors.New("pulse: AppID must be a non-empty string")
	// ErrMissingAppSecret indicates the AppSecret parameter was omitted.
	ErrMissingAppSecret = errors.New("pulse: AppSecret must be a non-empty string")
	// ErrSecretTooShort indicates the AppSecret does not meet minimum entropy requirements.
	ErrSecretTooShort = errors.New("pulse: AppSecret must be at least 16 characters in length")
	// ErrInvalidTimeout indicates the Timeout duration is non-positive.
	ErrInvalidTimeout = errors.New("pulse: Timeout must be greater than zero")
)

// Config defines connection and security settings required to instantiate the Pulse Go Server SDK.
// Strictly adheres to fail-fast principles with zero implicit default fallbacks.
type Config struct {
	BaseURL    string        // Base URL of the Pulse REST API (e.g. "http://localhost:32011")
	AppID      string        // Application / Tenant identifier
	AppSecret  string        // Cryptographic secret key used for JWT signing and webhook verification
	Timeout    time.Duration // Timeout for outgoing HTTP requests (must be > 0)
	MaxRetries int           // Maximum retry attempts for transient 5xx or rate limit errors
}

// TokenOptions specifies optional claims when minting a client JWT session token.
type TokenOptions struct {
	Role       string                 // User role: "user", "admin", "moderator"
	TTL        time.Duration          // Validity duration (e.g. 24 * time.Hour)
	Metadata   map[string]interface{} // Arbitrary tenant metadata claims
}

// User represents a public user profile on the Pulse platform.
type User struct {
	UserID      string                 `json:"userId"`
	DisplayName string                 `json:"displayName"`
	AvatarURL   string                 `json:"avatarUrl,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   string                 `json:"createdAt,omitempty"`
	UpdatedAt   string                 `json:"updatedAt,omitempty"`
	IsBanned    bool                   `json:"isBanned,omitempty"`
}

// UpsertUserRequest holds the input data for creating or updating a user identity.
type UpsertUserRequest struct {
	UserID      string                 `json:"userId"`
	DisplayName string                 `json:"displayName"`
	AvatarURL   string                 `json:"avatarUrl,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Conversation represents a direct or group conversation channel.
type Conversation struct {
	ConvID    string   `json:"convId"`
	Type      string   `json:"type"` // "DIRECT" or "GROUP"
	Name      string   `json:"name,omitempty"`
	ConvSeq   int64    `json:"convSeq"`
	MemberIDs []string `json:"memberIds"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

// CreateConversationRequest specifies parameters to create a new conversation.
type CreateConversationRequest struct {
	Type      string                 `json:"type"` // "DIRECT" or "GROUP"
	MemberIDs []string               `json:"memberIds"`
	Name      string                 `json:"name,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// SystemMessageRequest holds parameters for publishing a server-originated system message.
type SystemMessageRequest struct {
	ConvID   string                 `json:"convId"`
	Body     string                 `json:"body"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// BroadcastRequest holds parameters for cluster-wide or targeted broadcast alerts.
type BroadcastRequest struct {
	Body          string                 `json:"body"`
	TargetUserIDs []string               `json:"targetUserIds,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// MessageResponse represents the server confirmation for a delivered message.
type MessageResponse struct {
	MsgID     string `json:"msgId"`
	ConvID    string `json:"convId"`
	ConvSeq   int64  `json:"convSeq"`
	SenderID  string `json:"senderId"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
}

// WebhookEvent represents an inbound real-time event envelope emitted by Pulse.
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	AppID     string                 `json:"appId"`
	Payload   map[string]interface{} `json:"payload"`
}

