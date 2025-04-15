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
	clients := make(map[string]*client)
	return &messageServer{
		userRepo: userRepo,
		clients:  clients,
	}
}

func (s *messageServer) Start(ctx context.Context, stop context.CancelFunc) {

	go s.sendMessageProcess()

	// ticker := time.NewTicker(time.Second * 5)
	// done := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-done:
	// 			return
	// 		case <-ticker.C:
	// 			fmt.Println("DEBUG ticker", s.clients)
	// 			for id := range s.clients {
	// 				fmt.Println(id)
	// 			}
	// 		}
	// 	}
	// }()
	// <-done

	<-ctx.Done()

	log.Println("shutting down message server...")
}

func (s *messageServer) receiveMessageProcess(client *client) {
	defer func() {
		client.wg.Done()
		client.isClosed = true
	}()

	// first message is auth
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
			}
			return
		}
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
			// client.message <- message
		// comment for lazy fix lint
		// case websocket.PingMessage:
		// 	client.connection.WriteMessage(websocket.PongMessage, []byte{})
		default:
			log.Printf("websocket message received of type %d\n, ignored", messageType)
		}
	}
}

func (s *messageServer) sendMessageProcess() {
	for {
		for id, client := range s.clients {
			select {
			case msg, ok := <-client.message:
				if !ok {
					s.removeClient(id)
					continue
				}

				client.mu.Lock()
				if client.isClosed {
					client.mu.Unlock()
					s.removeClient(id)
					continue
				}

				if err := client.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
					client.isClosed = true
					log.Println("write error:", err)
					if err := client.connection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
						log.Print("write close err:", err)
					}
					client.connection.Close()
				}
				client.mu.Unlock()
			default:
			}

		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *messageServer) removeClient(id string) {
	s.wrmu.Lock()
	delete(s.clients, id)
	s.wrmu.Unlock()
}

func (s *messageServer) validateGoogleToken(token string) (domain.Profile, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("failed to create request to Google OAuth: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("failed to get profile from Google OAuth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Profile{}, fmt.Errorf("non-200 response from Google OAuth: %d", resp.StatusCode)
	}

	var profile domain.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return domain.Profile{}, fmt.Errorf("failed to decode profile from Google OAuth: %w", err)
	}

	// Ensure we have an email
	if profile.Email == "" {
		return domain.Profile{}, fmt.Errorf("profile doesn't contain email")
	}

	return profile, nil
}

func (s *messageServer) auth(client *client) error {
	messageType, message, err := client.connection.ReadMessage()
	if err != nil {
		log.Printf("authentication error: %v\n", err)
		return err
	}

	if messageType != websocket.TextMessage {
		log.Printf("expected text message for authentication, got type %d\n", messageType)
		return err
	}

	var authRequest dto.AuthRequest
	if err := json.Unmarshal(message, &authRequest); err != nil {
		log.Printf("invalid JSON auth format: %v\n", err)
		return err
	}

	profile, err := s.validateGoogleToken(authRequest.Token)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByProvider("GOOGLE", profile.Sub)
	if err != nil {
		return err
	}

	// Set client information
	client.userID = user.ID
	client.profile = profile

	// Register client in the map
	s.wrmu.Lock()
	s.clients[user.ID] = client
	s.wrmu.Unlock()

	// Send authentication success response
	err = client.connection.WriteMessage(websocket.TextMessage, []byte(profile.Sub))
	if err != nil {
		return err
	}
	return nil
}
