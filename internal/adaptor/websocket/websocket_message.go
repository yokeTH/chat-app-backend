package websocket

import "github.com/goccy/go-json"

type EventType string

const (
	EventTypeConnect        EventType = "connect"
	EventTypeMessage        EventType = "message"
	EventTypeDisconnect     EventType = "disconnect"
	EventTypeReactionAdd    EventType = "reaction_add"
	EventTypeReactionRemove EventType = "reaction_remove"
	EventTypeReadReceipt    EventType = "read_receipt"
	EventTypeTypingEnd      EventType = "typing_end"
	EventTypeTypingStart    EventType = "typing_start"
	EventTypeUserStatus     EventType = "user_status"
)

type WebSocketMessage struct {
	Event     EventType       `json:"event"`
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
	MessageType    string       `json:"type"`
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

type UserStatusType string

const (
	UserStatusTypeOnline  UserStatusType = "online"
	UserStatusTypeOffline UserStatusType = "offline"
)

type UserStatus struct {
	UserID string         `json:"userId,omitempty"`
	Status UserStatusType `json:"status"`
	Name   string         `json:"name,omitempty"`
}
