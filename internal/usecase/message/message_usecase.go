package message

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type messageUseCase struct {
	server *messageServer
}

func NewMessageUseCase(server *messageServer) *messageUseCase {
	return &messageUseCase{
		server: server,
	}
}

func (m *messageUseCase) RegisterClient(c *websocket.Conn) *sync.WaitGroup {
	m.server.wrmu.Lock()
	defer m.server.wrmu.Unlock()
	var wg sync.WaitGroup
	wg.Add(1)
	client := client{
		message:    make(chan []byte, 10),
		connection: c,
		wg:         &wg,
	}
	go m.server.receiveMessageProcess(&client)
	return &wg
}

func (m *messageUseCase) SendMessage(id string, message string) {
	client, ok := m.server.clients[id]
	if ok {
		client.message <- []byte(message)
		return
	}
	fmt.Printf("user %s does not exist\n", id)
}
