package websocket

import (
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func (s *messageServer) sendMessageProcess() {
	sendTicker := time.Tick(10 * time.Millisecond)
	pingTicker := time.Tick(1 * time.Minute)

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
			s.removeClientByID(id)

		case msg, ok := <-c.message:
			c.mu.Lock()

			if !ok || c.isClosed {
				c.isClosed = true
				c.mu.Unlock()
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("write error:", err)
				c.close()
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

		if err := c.connection.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
			log.Printf("Ping failed to user %s: %v", id, err)
			c.close()
		}

		c.mu.Unlock()
	}
}
