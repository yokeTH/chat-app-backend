package handler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/gofiber-template/internal/usecase/message"
)

type messageHandler struct {
	msgUseCase message.MessageUseCase
}

func NewMessageHandler(msgUseCase message.MessageUseCase) *messageHandler {
	return &messageHandler{
		msgUseCase: msgUseCase,
	}
}

func (h *messageHandler) HandleMessage(c *websocket.Conn) {
	h.msgUseCase.RegisterClient(c).Wait()
}
