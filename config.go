package pulse

import (
	"strings"
)

// Validate checks whether the provided Config satisfies all mandatory fail-fast invariants.
// Returns a descriptive error on any violation, strictly prohibiting silent fallback values.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}

	trimmedURL := strings.TrimSpace(c.BaseURL)
	if trimmedURL == "" || (!strings.HasPrefix(trimmedURL, "http://") && !strings.HasPrefix(trimmedURL, "https://")) {
		return ErrMissingBaseURL
	}

	if strings.TrimSpace(c.AppID) == "" {
		return ErrMissingAppID
	}

	if strings.TrimSpace(c.AppSecret) == "" {
		return ErrMissingAppSecret
	}

	if len(c.AppSecret) < 16 {
		return ErrSecretTooShort
	}

	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	return nil
}

