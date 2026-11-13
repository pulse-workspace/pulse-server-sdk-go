package pulse

import (
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr error
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: ErrNilConfig,
		},
		{
			name: "missing BaseURL",
			cfg: &Config{
				BaseURL:   "",
				AppID:     "app_1",
				AppSecret: "1234567890123456",
				Timeout:   5 * time.Second,
			},
			wantErr: ErrMissingBaseURL,
		},
		{
			name: "invalid BaseURL scheme",
			cfg: &Config{
				BaseURL:   "ftp://localhost:32011",
				AppID:     "app_1",
				AppSecret: "1234567890123456",
				Timeout:   5 * time.Second,
			},
			wantErr: ErrMissingBaseURL,
		},
		{
			name: "missing AppID",
			cfg: &Config{
				BaseURL:   "http://localhost:32011",
				AppID:     "   ",
				AppSecret: "1234567890123456",
				Timeout:   5 * time.Second,
			},
			wantErr: ErrMissingAppID,
		},
		{
			name: "missing AppSecret",
			cfg: &Config{
				BaseURL:   "http://localhost:32011",
				AppID:     "app_1",
				AppSecret: "",
				Timeout:   5 * time.Second,
			},
			wantErr: ErrMissingAppSecret,
		},
		{
			name: "AppSecret too short",
			cfg: &Config{
				BaseURL:   "http://localhost:32011",
				AppID:     "app_1",
				AppSecret: "short_secret",
				Timeout:   5 * time.Second,
			},
			wantErr: ErrSecretTooShort,
		},
		{
			name: "invalid zero timeout",
			cfg: &Config{
				BaseURL:   "http://localhost:32011",
				AppID:     "app_1",
				AppSecret: "1234567890123456",
				Timeout:   0,
			},
			wantErr: ErrInvalidTimeout,
		},
		{
			name: "valid configuration",
			cfg: &Config{
				BaseURL:   "http://localhost:32011",
				AppID:     "app_production",
				AppSecret: "0123456789abcdef0123456789abcdef",
				Timeout:   5 * time.Second,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

