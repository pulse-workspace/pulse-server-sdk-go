package pulse

import (
	"context"
	"errors"
	"strings"
)

var (
	// ErrInvalidConversationType indicates an unsupported conversation type was provided.
	ErrInvalidConversationType = errors.New("pulse: conversation type must be 'DIRECT' or 'GROUP'")
	// ErrEmptyMemberIDs indicates no participant IDs were supplied.
	ErrEmptyMemberIDs = errors.New("pulse: memberIds must contain at least one user ID")
	// ErrEmptyConvID indicates an empty conversation ID was passed.
	ErrEmptyConvID = errors.New("pulse: convID must be a non-empty string")
)

// ConversationService handles direct and group conversation lifecycle operations.
type ConversationService struct {
	http *httpClient
}

// newConversationService creates an instance of ConversationService.
func newConversationService(http *httpClient) *ConversationService {
	return &ConversationService{http: http}
}

// Create establishes a new direct (1-1) or group conversation channel.
func (s *ConversationService) Create(ctx context.Context, req CreateConversationRequest) (*Conversation, error) {
	if req.Type != "DIRECT" && req.Type != "GROUP" {
		return nil, ErrInvalidConversationType
	}
	if len(req.MemberIDs) == 0 {
		return nil, ErrEmptyMemberIDs
	}

	var conv Conversation
	err := s.http.request(ctx, "POST", "/api/v1/conversations", req, &conv)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// Get fetches details and member roster of an existing conversation.
func (s *ConversationService) Get(ctx context.Context, convID string) (*Conversation, error) {
	if strings.TrimSpace(convID) == "" {
		return nil, ErrEmptyConvID
	}

	var conv Conversation
	err := s.http.request(ctx, "GET", "/api/v1/conversations/"+escapePath(convID), nil, &conv)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// AddMembers appends new user IDs to an existing group conversation.
func (s *ConversationService) AddMembers(ctx context.Context, convID string, memberIDs []string) error {
	if strings.TrimSpace(convID) == "" {
		return ErrEmptyConvID
	}
	if len(memberIDs) == 0 {
		return ErrEmptyMemberIDs
	}

	body := map[string][]string{"memberIds": memberIDs}
	return s.http.request(ctx, "POST", "/api/v1/conversations/"+escapePath(convID)+"/members", body, nil)
}

// RemoveMembers removes participants from a group conversation.
func (s *ConversationService) RemoveMembers(ctx context.Context, convID string, memberIDs []string) error {
	if strings.TrimSpace(convID) == "" {
		return ErrEmptyConvID
	}
	if len(memberIDs) == 0 {
		return ErrEmptyMemberIDs
	}

	body := map[string][]string{"memberIds": memberIDs}
	return s.http.request(ctx, "DELETE", "/api/v1/conversations/"+escapePath(convID)+"/members", body, nil)
}

// Archive marks a conversation as archived, making it immutable.
func (s *ConversationService) Archive(ctx context.Context, convID string) error {
	if strings.TrimSpace(convID) == "" {
		return ErrEmptyConvID
	}

	return s.http.request(ctx, "POST", "/api/v1/conversations/"+escapePath(convID)+"/archive", nil, nil)
}

