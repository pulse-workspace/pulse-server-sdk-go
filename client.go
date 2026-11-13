package pulse

// Client provides the primary entrypoint for the Pulse Go Server SDK.
// Encapsulates authentication, identity profiles, conversation orchestration, and messaging.
type Client struct {
	config        Config
	Auth          *AuthService
	Users         *UserService
	Conversations *ConversationService
	Messages      *MessageService
	http          *httpClient
}

// NewClient instantiates a new Client after validating fail-fast configuration invariants.
// Returns an error if any required setting is missing or invalid.
func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	httpTransport := newHTTPClient(cfg)

	return &Client{
		config:        cfg,
		Auth:          NewAuthService(cfg.AppID, cfg.AppSecret),
		Users:         newUserService(httpTransport),
		Conversations: newConversationService(httpTransport),
		Messages:      newMessageService(httpTransport),
		http:          httpTransport,
	}, nil
}

// VerifyWebhook validates an inbound webhook payload signature against the client's configured secret.
func (c *Client) VerifyWebhook(payload []byte, signature string) bool {
	return VerifyWebhookSignature(payload, signature, c.config.AppSecret)
}

// ComputeWebhookSignature calculates the expected signature for an outbound webhook payload.
func (c *Client) ComputeWebhookSignature(payload []byte) string {
	return ComputeWebhookSignature(payload, c.config.AppSecret)
}

