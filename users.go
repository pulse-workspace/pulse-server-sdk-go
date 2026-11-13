package pulse

import (
	"context"
	"errors"
	"strings"
)

var (
	// ErrMissingDisplayName indicates displayName was empty when upserting a user.
	ErrMissingDisplayName = errors.New("pulse: displayName is required when upserting a user")
)

// UserService manages user profiles, identity synchronizations, and moderation actions.
type UserService struct {
	http *httpClient
}

// newUserService instantiates a new UserService with the provided HTTP transport.
func newUserService(http *httpClient) *UserService {
	return &UserService{http: http}
}

// Upsert creates or updates a user profile on the Pulse platform.
func (s *UserService) Upsert(ctx context.Context, req UpsertUserRequest) (*User, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return nil, ErrEmptyUserID
	}
	if strings.TrimSpace(req.DisplayName) == "" {
		return nil, ErrMissingDisplayName
	}

	var user User
	err := s.http.request(ctx, "PUT", "/api/v1/users", req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Get fetches the public profile and current presence status for a given user ID.
func (s *UserService) Get(ctx context.Context, userID string) (*User, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrEmptyUserID
	}

	var user User
	err := s.http.request(ctx, "GET", "/api/v1/users/"+escapePath(userID), nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Ban restricts a user from connecting or dispatching messages.
func (s *UserService) Ban(ctx context.Context, userID, reason string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrEmptyUserID
	}
	if reason == "" {
		reason = "Administrative ban"
	}

	body := map[string]string{"reason": reason}
	return s.http.request(ctx, "POST", "/api/v1/users/"+escapePath(userID)+"/ban", body, nil)
}

// Unban restores access for a previously restricted user account.
func (s *UserService) Unban(ctx context.Context, userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrEmptyUserID
	}

	return s.http.request(ctx, "POST", "/api/v1/users/"+escapePath(userID)+"/unban", nil, nil)
}

// RevokeSessions terminates all active WebSocket sessions for a user across all cluster gateways.
func (s *UserService) RevokeSessions(ctx context.Context, userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrEmptyUserID
	}

	return s.http.request(ctx, "POST", "/api/v1/users/"+escapePath(userID)+"/revoke-sessions", nil, nil)
}

