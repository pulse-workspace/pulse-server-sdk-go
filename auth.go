package pulse

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrEmptyUserID indicates an empty user ID was passed to token creation.
	ErrEmptyUserID = errors.New("pulse: userID must be a non-empty string")
	// ErrInvalidTokenFormat indicates a malformed JWT string structure.
	ErrInvalidTokenFormat = errors.New("pulse: malformed JWT token: must contain exactly 3 dot-separated segments")
	// ErrInvalidSignature indicates a cryptographic signature mismatch.
	ErrInvalidSignature = errors.New("pulse: invalid JWT token signature")
	// ErrTokenExpired indicates the token expiration timestamp is in the past.
	ErrTokenExpired = errors.New("pulse: JWT token has expired")
)

// AuthService handles creation and cryptographic verification of client JWT session tokens.
type AuthService struct {
	appID     string
	appSecret string
}

// NewAuthService instantiates an AuthService with the tenant credentials.
func NewAuthService(appID, appSecret string) *AuthService {
	return &AuthService{
		appID:     appID,
		appSecret: appSecret,
	}
}

// CreateToken mints an HMAC-SHA256 signed JWT token for a specific user ID.
func (a *AuthService) CreateToken(userID string, opts ...TokenOptions) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", ErrEmptyUserID
	}

	role := "user"
	ttl := 24 * time.Hour
	var metadata map[string]interface{}

	if len(opts) > 0 {
		opt := opts[0]
		if opt.Role != "" {
			role = opt.Role
		}
		if opt.TTL > 0 {
			ttl = opt.TTL
		}
		if opt.Metadata != nil {
			metadata = opt.Metadata
		}
	}

	now := time.Now().Unix()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload := map[string]interface{}{
		"sub":  userID,
		"iss":  a.appID,
		"aud":  "pulse-client",
		"role": role,
		"iat":  now,
		"exp":  now + int64(ttl.Seconds()),
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("pulse: failed to marshal JWT header: %w", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("pulse: failed to marshal JWT payload: %w", err)
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signInput := headerB64 + "." + payloadB64

	h := hmac.New(sha256.New, []byte(a.appSecret))
	h.Write([]byte(signInput))
	sigB64 := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return signInput + "." + sigB64, nil
}

// VerifyToken validates the signature and expiration timestamp of a token, returning its parsed claims.
func (a *AuthService) VerifyToken(tokenStr string) (map[string]interface{}, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}

	signInput := parts[0] + "." + parts[1]
	expectedSig := computeHMAC(signInput, a.appSecret)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, ErrInvalidSignature
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("pulse: failed to decode JWT payload: %w", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("pulse: failed to unmarshal JWT claims: %w", err)
	}

	if expVal, ok := claims["exp"]; ok {
		if expFloat, ok := expVal.(float64); ok {
			if int64(expFloat) < time.Now().Unix() {
				return nil, ErrTokenExpired
			}
		}
	}

	return claims, nil
}

// computeHMAC computes a base64url-encoded HMAC-SHA256 digest.
func computeHMAC(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

