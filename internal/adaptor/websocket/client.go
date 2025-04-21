package websocket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/chat-app-backend/internal/domain"
)

type client struct {
	id         string
	isClosed   bool
	terminate  chan bool
	wg         *sync.WaitGroup
	mu         sync.Mutex
	connection *websocket.Conn
	message    chan []byte
	userID     string
	profile    domain.Profile
}

func (c *client) sendError(message string) {
	_ = c.connection.WriteMessage(websocket.TextMessage, []byte(message))
	c.connection.Close()
}

func (c *client) close() {
	_ = c.connection.WriteMessage(websocket.CloseMessage, []byte{})
	c.terminate <- true
	c.isClosed = true
	c.connection.Close()
}
