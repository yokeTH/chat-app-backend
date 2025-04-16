package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/user"
)

type messageServer struct {
	userRepo user.UserRepository
	clients  map[string]*client
	wrmu     sync.RWMutex
}

func NewMessageServer(userRepo user.UserRepository) *messageServer {
	return &messageServer{
		userRepo: userRepo,
		clients:  make(map[string]*client),
	}
}

func (s *messageServer) Start(ctx context.Context, stop context.CancelFunc) {
	go s.sendMessageProcess()

	<-ctx.Done()
	log.Println("shutting down message server...")
}

func (s *messageServer) receiveMessageProcess(client *client) {
	defer func() {
		client.wg.Done()
		client.isClosed = true
	}()

	// First message must be auth
	if err := s.auth(client); err != nil {
		log.Printf("Authentication failed: %v", err)
		client.sendError("Authentication failed")
		return
	}

	_ = client.connection.SetReadDeadline(time.Now().Add(35 * time.Second))

	client.connection.SetPongHandler(func(appData string) error {
		log.Printf("Received pong from user: %s", client.userID)
		return client.connection.SetReadDeadline(time.Now().Add(35 * time.Second))
	})

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

		_ = client.connection.SetReadDeadline(time.Now().Add(35 * time.Second))

		switch messageType {
		case websocket.TextMessage:
			msg := string(message)
			msgs := strings.Split(msg, ":")
			if len(msgs) < 2 {
				log.Printf("malformed message received: %s\n", msg)
				continue
			}
			id := msgs[0]
			text := strings.Join(msgs[1:], ":")
			fmt.Println(id, text)
			client.message <- []byte(text)
		default:
			log.Printf("websocket message type %d ignored\n", messageType)
		}
	}
}

func (s *messageServer) sendMessageProcess() {
	sendTicker := time.Tick(10 * time.Millisecond)
	pingTicker := time.Tick(10 * time.Second)

	for {
		select {
		case <-sendTicker:
			s.broadcastMessages()

		case <-pingTicker:
			s.sendPings()
		}
	}
}

func (s *messageServer) broadcastMessages() {
	for id, c := range s.clients {
		select {
		case <-c.terminate:
			s.removeClient(id)

		case msg, ok := <-c.message:
			c.mu.Lock()

			if !ok || c.isClosed {
				c.isClosed = true
				c.mu.Unlock()
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("write error:", err)
				c.terminate <- true
				_ = c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				c.connection.Close()
				c.isClosed = true
			}

			c.mu.Unlock()

		default:
		}
	}
}

func (s *messageServer) sendPings() {
	for id, c := range s.clients {
		c.mu.Lock()

		if c.isClosed {
			c.mu.Unlock()
			continue
		}

		c.connection.SetPongHandler(func(appData string) error {
			log.Printf("Received pong from user: %s", c.userID)
			return c.connection.SetReadDeadline(time.Now().Add(30 * time.Second))
		})

		if err := c.connection.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
			log.Printf("Failed to set read deadline for %s: %v", id, err)
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

	delete(s.clients, id)
	_ = s.userRepo.SetIsOnline(id, false)
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

	userData, err := s.userRepo.GetUserByProvider("GOOGLE", profile.Sub)
	if err != nil {
		return err
	}

	if err := s.userRepo.SetIsOnline(userData.ID, true); err != nil {
		return err
	}

	c.userID = userData.ID
	c.profile = profile

	s.wrmu.Lock()
	s.clients[userData.ID] = c
	s.wrmu.Unlock()

	return c.connection.WriteMessage(websocket.TextMessage, []byte(profile.Sub))
}

func (s *messageServer) validateGoogleToken(token string) (domain.Profile, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Profile{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var profile domain.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return domain.Profile{}, fmt.Errorf("decoding response: %w", err)
	}

	if profile.Email == "" {
		return domain.Profile{}, fmt.Errorf("missing email in profile")
	}

	return profile, nil
}
