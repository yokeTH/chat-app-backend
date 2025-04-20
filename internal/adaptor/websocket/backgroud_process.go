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
