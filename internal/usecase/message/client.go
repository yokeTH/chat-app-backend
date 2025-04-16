package message

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/gofiber-template/internal/domain"
)

type client struct {
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
	if err := c.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		c.isClosed = true
		log.Println("write error:", err)
		if err := c.connection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			log.Print("write close err:", err)
		}
		c.connection.Close()
	}
	c.isClosed = true
	c.connection.Close()
}
