package pulse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// httpClient manages robust HTTP communication with the Pulse cluster.
type httpClient struct {
	baseURL    string
	appID      string
	appSecret  string
	timeout    time.Duration
	maxRetries int
	client     *http.Client
}

// newHTTPClient creates an HTTP client using the supplied configuration.
func newHTTPClient(cfg Config) *httpClient {
	retries := cfg.MaxRetries
	if retries <= 0 {
		retries = 3
	}

	return &httpClient{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		appID:      cfg.AppID,
		appSecret:  cfg.AppSecret,
		timeout:    cfg.Timeout,
		maxRetries: retries,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// request executes an HTTP request with context timeout, authentication headers, and exponential retry.
func (c *httpClient) request(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	reqURL := c.baseURL + path

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("pulse: failed to marshal request body: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		reqCtx, cancel := context.WithTimeout(ctx, c.timeout)

		var bodyReader io.Reader
		if reqBody != nil {
			bodyReader = bytes.NewReader(reqBody)
		}

		req, err := http.NewRequestWithContext(reqCtx, method, reqURL, bodyReader)
		if err != nil {
			cancel()
			return fmt.Errorf("pulse: failed to create HTTP request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Pulse-App-Id", c.appID)
		req.Header.Set("Authorization", "Bearer "+c.appSecret)

		resp, err := c.client.Do(req)
		if err != nil {
			cancel()
			lastErr = err
			if attempt < c.maxRetries {
				sleepBackoff(attempt)
				continue
			}
			return fmt.Errorf("pulse: request failed after %d attempts: %w", attempt+1, lastErr)
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		cancel()

		if readErr != nil {
			lastErr = readErr
			if attempt < c.maxRetries {
				sleepBackoff(attempt)
				continue
			}
			return fmt.Errorf("pulse: failed to read response body: %w", readErr)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("pulse: API error [status %d]: %s", resp.StatusCode, string(respBody))
			// Retry only on 5xx or 429 rate limit errors
			if (resp.StatusCode >= 500 || resp.StatusCode == 429) && attempt < c.maxRetries {
				sleepBackoff(attempt)
				continue
			}
			return lastErr
		}

		if out != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, out); err != nil {
				return fmt.Errorf("pulse: failed to unmarshal response: %w", err)
			}
		}

		return nil
	}

	return lastErr
}

// sleepBackoff implements exponential backoff delay based on retry attempt index.
func sleepBackoff(attempt int) {
	delay := time.Duration(math.Pow(2, float64(attempt))*100) * time.Millisecond
	time.Sleep(delay)
}

// escapePath safely escapes URL path segments.
func escapePath(seg string) string {
	return url.PathEscape(seg)
}

