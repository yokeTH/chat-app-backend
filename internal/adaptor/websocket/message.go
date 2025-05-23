package websocket

import (
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
)

func (s *messageServer) handleEventTypeMessage(payload json.RawMessage) error {
	var chatMsg ChatMessage
	if err := json.Unmarshal(payload, &chatMsg); err != nil {
		log.Printf("invalid chat message payload: %v", err)
		return err
	}

	log.Printf("received message from %s: %s", chatMsg.SenderID, chatMsg.Content)
	content := dto.CreateMessageRequest{
		ConversationID: chatMsg.ConversationID,
		Content:        chatMsg.Content,
	}

	createdMessage, err := s.messageUC.Create(chatMsg.SenderID, content)
	if err != nil {
		log.Printf("failed to create message: %v", err)
		return err
	}

	createdMessageResponse, err := s.messageDto.ToResponse(createdMessage)
	if err != nil {
		log.Printf("failed to transform to dto: %v", err)
		return err
	}

	payloadResponse, err := json.Marshal(createdMessageResponse)
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return err
	}

	createdMessageJson, err := json.Marshal(WebSocketMessage{
		Event:     EventTypeMessage,
		Payload:   payloadResponse,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return err
	}

	if err := s.BroadcastToMembersInConversation(chatMsg.ConversationID, createdMessageJson); err != nil {
		return err
	}

	return nil
}

func (s *messageServer) BroadcastToMembersInConversation(conversationID string, msg []byte) error {
	members, err := s.conversationUC.GetMembers(conversationID)
	if err != nil {
		log.Printf("failed to get conversation members : %v", err)
		return err
	}

	for _, member := range *members {
		s.sendMessageToUserID(member.ID, msg)
	}
	return nil
}

func (s *messageServer) handleEventTypeTyping(payload json.RawMessage, currentUserID string, isTyping bool) error {
	var typing TypingEvent
	if err := json.Unmarshal(payload, &typing); err != nil {
		log.Printf("invalid typing_start payload: %v", err)
		return err
	}

	members, err := s.conversationUC.GetMembers(typing.ConversationID)
	if err != nil {
		log.Printf("failed to get conversation members : %v", err)
		return err
	}

	var event EventType
	if isTyping {
		event = EventTypeTypingStart
	} else {
		event = EventTypeTypingEnd
	}
	msg, err := json.Marshal(WebSocketMessage{
		Event:     event,
		Payload:   payload,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return err
	}

	for _, member := range *members {
		if member.ID != currentUserID {
			s.sendMessageToUserID(member.ID, msg)
		}
	}

	if isTyping {
		log.Printf("user %s started typing in conversation %s", typing.UserID, typing.ConversationID)
	} else {
		log.Printf("user %s ended typing in conversation %s", typing.UserID, typing.ConversationID)
	}
	return nil
}

func (s *messageServer) broadcastUserStatus(userID string, isOnline bool) {
	var status UserStatusType
	if isOnline {
		status = UserStatusTypeOnline
	} else {
		status = UserStatusTypeOffline
	}

	userStatusString, err := json.Marshal(UserStatus{
		UserID: userID,
		Status: status,
	})
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return
	}

	respMsg, err := json.Marshal(WebSocketMessage{
		Event:     EventTypeUserStatus,
		Payload:   userStatusString,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return
	}

	s.broadcast(respMsg)
}

func (s *messageServer) broadcast(message []byte) {
	for _, client := range s.clients {
		client.mu.Lock()
		client.message <- message
		client.mu.Unlock()
	}
}

func (s *messageServer) BoardcastConversation(conversation dto.ConversationResponse) {
	payloadResponse, err := json.Marshal(conversation)
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return
	}

	wsMsg, err := json.Marshal(WebSocketMessage{
		Event:     EventTypeConversationUpdate,
		Payload:   payloadResponse,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return
	}
	s.broadcast(wsMsg)
}

func (s *messageServer) BroadcastName(userID, name string) {
	payload, err := json.Marshal(UserStatus{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		return
	}

	msg, err := json.Marshal(WebSocketMessage{
		Event:     EventTypeUserStatus,
		Payload:   payload,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		return
	}
	s.broadcast(msg)
}
