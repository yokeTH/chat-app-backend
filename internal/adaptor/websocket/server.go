package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/conversation"
	"github.com/yokeTH/gofiber-template/internal/usecase/message"
	"github.com/yokeTH/gofiber-template/internal/usecase/user"
)

type messageServer struct {
	userUC         user.UserUseCase
	messageUC      message.MessageUseCase
	conversationUC conversation.ConversationUseCase
	messageDto     dto.MessageDto
	clients        map[string]*Client
	wrmu           sync.RWMutex
}

func NewMessageServer(userUC user.UserUseCase, messageUC message.MessageUseCase, conversationUC conversation.ConversationUseCase, messageDto dto.MessageDto) *messageServer {
	return &messageServer{
		userUC:         userUC,
		messageUC:      messageUC,
		conversationUC: conversationUC,
		messageDto:     messageDto,
		clients:        make(map[string]*Client),
	}
}

func (s *messageServer) Start(ctx context.Context, stop context.CancelFunc) {
	go s.sendMessageProcess()

	<-ctx.Done()
	log.Println("shutting down message server...")
}

func (s *messageServer) receiveMessageProcess(client *Client) {
	defer func() {
		client.isClosed = true
	}()

	// First message must be auth
	if err := s.auth(client); err != nil {
		log.Printf("Authentication failed: %v", err)
		client.sendError("Authentication failed")
		return
	}

	for {
		messageType, message, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message error: %v\n", err)
			} else {
				log.Printf("user %s connection closed: %v", client.userID, err)
			}
			client.isClosed = true
			client.terminate <- true
			return
		}

		switch messageType {
		case websocket.TextMessage:
			var wsMsg WebSocketMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				log.Printf("invalid WebSocket message: %s, error: %v\n", string(message), err)
				continue
			}

			switch wsMsg.Event {
			case "message":
				var chatMsg ChatMessage
				if err := json.Unmarshal(wsMsg.Payload, &chatMsg); err != nil {
					log.Printf("invalid chat message payload: %v", err)
					continue
				}
				log.Printf("received message from %s: %s", chatMsg.SenderID, chatMsg.Content)
				content := dto.CreateMessageRequest{
					ConversationID: chatMsg.ConversationID,
					Content:        chatMsg.Content,
				}

				createdMessage, err := s.messageUC.Create(chatMsg.SenderID, content)
				if err != nil {
					log.Printf("failed to create message: %v", err)
					continue
				}
				createdMessageResponse, _ := s.messageDto.ToResponse(createdMessage)
				payload, _ := json.Marshal(createdMessageResponse)
				createdMessageJson, _ := json.Marshal(WebSocketMessage{
					Event:     "message",
					Payload:   payload,
					CreatedAt: time.Now().UnixMilli(),
				})
				members, _ := s.conversationUC.GetMembers(chatMsg.ConversationID)
				for _, member := range *members {
					if client := s.getClient(member.ID); client != nil {
						client.message <- createdMessageJson
					}
				}
				// client.message <- createdMessageJson
			case "typing_start":
				var typing TypingEvent
				if err := json.Unmarshal(wsMsg.Payload, &typing); err != nil {
					log.Printf("invalid typing_start payload: %v", err)
					continue
				}
				log.Printf("user %s started typing in conversation %s", typing.UserID, typing.ConversationID)
			default:
				log.Printf("unhandled WebSocket event: %s", wsMsg.Event)
			}

		default:
			log.Printf("websocket message type %d ignored\n", messageType)
		}
	}
}

func (m *messageServer) RegisterClient(c *websocket.Conn) *sync.WaitGroup {
	m.wrmu.Lock()
	defer m.wrmu.Unlock()
	var wg sync.WaitGroup
	wg.Add(1)
	client := Client{
		message:    make(chan []byte, 10),
		connection: c,
		wg:         &wg,
		terminate:  make(chan bool, 1),
	}
	go m.receiveMessageProcess(&client)
	return &wg
}

func (m *messageServer) HandleWebsocket(c *websocket.Conn) {
	m.RegisterClient(c).Wait()
}

func (m *messageServer) SendMessage(id string, message string) {
	client := m.getClient(id)
	if client != nil {
		client.message <- []byte(message)
	}
}

func (m *messageServer) getClient(id string) *Client {
	client, ok := m.clients[id]
	if ok {
		return client
	}
	return nil
}

func (s *messageServer) sendPings() {
	for id, c := range s.clients {
		c.mu.Lock()

		if c.isClosed {
			c.mu.Unlock()
			continue
		}

		if err := c.connection.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
			log.Printf("Ping failed to user %s: %v", id, err)
			c.isClosed = true
			c.terminate <- true
			_ = c.connection.WriteMessage(websocket.CloseMessage, []byte{})
			_ = c.connection.Close()
		}

		c.mu.Unlock()
	}
}

func (s *messageServer) removeClient(id string) {
	s.wrmu.Lock()
	defer s.wrmu.Unlock()
	client := s.getClient(id)
	if client != nil {
		client.wg.Done()
		delete(s.clients, id)
		_ = s.userUC.SetUserOffline(id)
	}
}

func (s *messageServer) auth(c *Client) error {
	msgType, data, err := c.connection.ReadMessage()
	if err != nil {
		log.Printf("auth read error: %v", err)
		return err
	}

	if msgType != websocket.TextMessage {
		return fmt.Errorf("invalid message type: %d", msgType)
	}

	var auth dto.AuthRequest
	if err := json.Unmarshal(data, &auth); err != nil {
		log.Printf("auth unmarshal error: %v", err)
		return err
	}

	profile, err := s.validateGoogleToken(auth.Token)
	if err != nil {
		return err
	}

	userData, err := s.userUC.GetGoogleProfile(profile.Sub)
	if err != nil {
		return err
	}

	if err := s.userUC.SetUserOnline(userData.ID); err != nil {
		return err
	}

	c.userID = userData.ID
	c.profile = *profile

	s.wrmu.Lock()
	s.clients[userData.ID] = c
	s.wrmu.Unlock()

	return c.connection.WriteMessage(websocket.TextMessage, []byte(profile.Sub))
}

func (s *messageServer) validateGoogleToken(token string) (*domain.Profile, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var profile domain.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if profile.Email == "" {
		return nil, fmt.Errorf("missing email in profile")
	}

	return &profile, nil
}
