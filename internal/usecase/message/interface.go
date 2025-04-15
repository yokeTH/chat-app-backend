package message

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type MessageRepository interface {
}

type MessageUseCase interface {
	RegisterClient(id string, c *websocket.Conn) *sync.WaitGroup
	SendMessage(id string, message string)
}
