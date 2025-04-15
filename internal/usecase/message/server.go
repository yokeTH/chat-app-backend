package message

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type client struct {
	isClosed   bool
	wg         *sync.WaitGroup
	mu         sync.Mutex
	connection *websocket.Conn
	message    chan []byte
}

type messageServer struct {
	clients map[string]*client
	wrmu    sync.RWMutex
}

func NewMessageServer() *messageServer {
	clients := make(map[string]*client)
	return &messageServer{
		clients: clients,
	}
}

func (s *messageServer) Start(ctx context.Context, stop context.CancelFunc) {

	go s.sendMessageProcess()

	<-ctx.Done()

	log.Println("shutting down message server...")
}

func (s *messageServer) ReceiveMessageProcess(client *client) {
	defer func() {
		client.wg.Done()
		client.isClosed = true
	}()

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
			client.message <- message
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
			msg, ok := <-client.message
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

		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *messageServer) removeClient(id string) {
	s.wrmu.Lock()
	delete(s.clients, id)
	s.wrmu.Unlock()
}
