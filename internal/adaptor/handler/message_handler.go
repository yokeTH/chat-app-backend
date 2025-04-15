package handler

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type messageHandler struct {
}

func NewMessageHandler() *messageHandler {
	return &messageHandler{}
}

func (h *messageHandler) HandleMessage(c *websocket.Conn) {
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message error: %v\n", err)
			}
			return
		}
		if messageType == websocket.TextMessage {
			// Broadcast the received message
			// broadcast <- string(message)
			fmt.Println(string(message))
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
