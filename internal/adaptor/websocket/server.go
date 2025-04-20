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
	clients        map[string]*client
	wrmu           sync.RWMutex
}

func NewMessageServer(userUC user.UserUseCase, messageUC message.MessageUseCase, conversationUC conversation.ConversationUseCase, messageDto dto.MessageDto) *messageServer {
	return &messageServer{
		userUC:         userUC,
		messageUC:      messageUC,
		conversationUC: conversationUC,
		messageDto:     messageDto,
		clients:        make(map[string]*client),
	}
}

func (s *messageServer) Start(ctx context.Context, stop context.CancelFunc) {
	go s.sendMessageProcess()

	<-ctx.Done()
	log.Println("shutting down message server...")
}

func (s *messageServer) receiveMessageProcess(uuid string, client *client) {
	defer func() {
		client.isClosed = true
	}()

	// First message must be auth
	if err := s.auth(client); err != nil {
		log.Printf("Authentication failed: %v", err)
		client.sendError("Authentication failed")
		return
	}

	s.addClient(uuid, client)

	for {
		messageType, message, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message error: %v\n", err)
			} else {
				log.Printf("user %s connection closed: %v", client.userID, err)
			}
			client.close()
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
			case EventTypeMessage:
				if err := s.handleEventTypeMessage(wsMsg.Payload); err != nil {
					continue
				}
			case EventTypeTypingStart:
				if err := s.handleEventTypeTyping(wsMsg.Payload, client.userID, false); err != nil {
					continue
				}
			case EventTypeTypingEnd:
				if err := s.handleEventTypeTyping(wsMsg.Payload, client.userID, false); err != nil {
					continue
				}
			default:
				log.Printf("unhandled WebSocket event: %s", wsMsg.Event)
			}

		default:
			log.Printf("websocket message type %d ignored\n", messageType)
		}
	}
}

func (m *messageServer) registerClient(c *websocket.Conn) *sync.WaitGroup {
	m.wrmu.Lock()
	defer m.wrmu.Unlock()
	var wg sync.WaitGroup
	wg.Add(1)
	requestid := c.Locals("requestid").(string)
	client := client{
		id:         requestid,
		message:    make(chan []byte, 10),
		connection: c,
		wg:         &wg,
		terminate:  make(chan bool, 1),
	}
	go m.receiveMessageProcess(requestid, &client)
	return &wg
}

func (m *messageServer) HandleWebsocket(c *websocket.Conn) {
	m.registerClient(c).Wait()
}

func (m *messageServer) sendMessageToUserID(id string, message []byte) {
	for _, client := range m.getClientByUserID(id) {
		if client != nil {
			client.message <- message
		}
	}
}

func (m *messageServer) getClientByUserID(userID string) []*client {
	var clients []*client
	for _, client := range m.clients {
		if client != nil && client.userID == userID {
			clients = append(clients, client)
		}
	}
	return clients
}

//nolint:unused
func (s *messageServer) removeClientByUserID(id string) {
	s.wrmu.Lock()
	defer s.wrmu.Unlock()
	clients := s.getClientByUserID(id)
	for _, client := range clients {
		if client != nil {
			client.wg.Done()
			delete(s.clients, client.id)
			_ = s.userUC.SetUserOffline(id)
		}
	}
}

func (s *messageServer) removeClientByID(id string) {
	s.wrmu.Lock()
	defer s.wrmu.Unlock()
	client, ok := s.clients[id]
	if !ok {
		return
	}
	client.wg.Done()
	delete(s.clients, id)
	userID := client.userID
	cnt := 0
	for _, client := range s.clients {
		if client.userID == userID {
			cnt++
		}
	}
	if cnt == 0 {
		_ = s.userUC.SetUserOffline(userID)
		go s.boardcaseUserStatus(client.userID, false)
	}
}

func (s *messageServer) auth(c *client) error {
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

	return nil
}

func (s *messageServer) addClient(uuid string, client *client) {
	go s.boardcaseUserStatus(client.userID, true)
	s.wrmu.Lock()
	s.clients[uuid] = client
	s.wrmu.Unlock()
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
