package pulse

import (
	"context"
	"errors"
	"strings"
)

var (
	// ErrEmptyBody indicates an empty message body was passed.
	ErrEmptyBody = errors.New("pulse: message body cannot be empty")
	// ErrEmptyMsgID indicates an empty message ID was provided.
	ErrEmptyMsgID = errors.New("pulse: msgID must be a non-empty string")
)

// MessageService handles server-originated messaging, broadcast pushes, and message retractions.
type MessageService struct {
	http *httpClient
}

// newMessageService instantiates MessageService with an active HTTP client.
func newMessageService(http *httpClient) *MessageService {
	return &MessageService{http: http}
}

// SendSystemMessage publishes an official server-originated system message into a conversation.
func (s *MessageService) SendSystemMessage(ctx context.Context, req SystemMessageRequest) (*MessageResponse, error) {
	if strings.TrimSpace(req.ConvID) == "" {
		return nil, ErrEmptyConvID
	}
	if strings.TrimSpace(req.Body) == "" {
		return nil, ErrEmptyBody
	}

	var resp MessageResponse
	err := s.http.request(ctx, "POST", "/api/v1/messages/system", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Broadcast sends an alert message to multiple users or the entire active cluster.
// Returns the number of successfully delivered broadcast messages.
func (s *MessageService) Broadcast(ctx context.Context, req BroadcastRequest) (int, error) {
	if strings.TrimSpace(req.Body) == "" {
		return 0, ErrEmptyBody
	}

	var resp struct {
		Success         bool `json:"success"`
		DispatchedCount int  `json:"dispatchedCount"`
	}
	err := s.http.request(ctx, "POST", "/api/v1/messages/broadcast", req, &resp)
	if err != nil {
		return 0, err
	}
	return resp.DispatchedCount, nil
}

// Recall revokes a previously dispatched message from a conversation for all members.
func (s *MessageService) Recall(ctx context.Context, convID, msgID, reason string) error {
	if strings.TrimSpace(convID) == "" {
		return ErrEmptyConvID
	}
	if strings.TrimSpace(msgID) == "" {
		return ErrEmptyMsgID
	}
	if reason == "" {
		reason = "Administrative recall"
	}

	body := map[string]string{"reason": reason}
	path := "/api/v1/conversations/" + escapePath(convID) + "/messages/" + escapePath(msgID) + "/recall"
	return s.http.request(ctx, "POST", path, body, nil)
}

