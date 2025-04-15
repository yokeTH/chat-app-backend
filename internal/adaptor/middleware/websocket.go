package middleware

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebsocketMiddleware interface {
}

type websocketMiddleware struct{}

func NewWebsocketMiddleware() *websocketMiddleware {
	return &websocketMiddleware{}
}

func (m *websocketMiddleware) RequiredUpgradeProtocal(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUpgradeRequired)
}
