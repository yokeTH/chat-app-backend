package websocket

import "github.com/goccy/go-json"

type WebSocketMessage struct {
	Event     string          `json:"event"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt int64           `json:"created_at"`
}

type ChatMessage struct {
	ID             string       `json:"id"`
	ConversationID string       `json:"conversationId"`
	Content        string       `json:"content"`
	SenderID       string       `json:"senderId"`
	Timestamp      int64        `json:"created_at"`
	Attachments    []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	Size     *int64 `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type TypingEvent struct {
	ConversationID string `json:"conversationId"`
	UserID         string `json:"userId"`
}
