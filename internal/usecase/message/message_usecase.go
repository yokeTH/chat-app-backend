package message

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type messageUseCase struct {
	server *messageServer
	repo   MessageRepository
}

func NewMessageUseCase(server *messageServer, repo MessageRepository) *messageUseCase {
	return &messageUseCase{
		server: server,
		repo:   repo,
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
		terminate:  make(chan bool, 1),
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

func (uc *messageUseCase) Create(senderID string, req dto.CreateMessageRequest) (*domain.Message, error) {
	message := &domain.Message{
		ConversationID: req.ConversationID,
		SenderID:       senderID,
		Content:        req.Content,
	}

	if err := uc.repo.Create(message); err != nil {
		return nil, apperror.InternalServerError(err, "failed to create message")
	}

	return message, nil
}

func (uc *messageUseCase) GetByID(id string) (*domain.Message, error) {
	message, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, apperror.NotFoundError(err, "message not found")
	}

	return message, nil
}

func (uc *messageUseCase) GetByConversationID(convoID string) (*[]domain.Message, error) {
	messages, err := uc.repo.FindByConversationID(convoID)
	if err != nil {
		return nil, apperror.InternalServerError(err, "failed to fetch messages")
	}
	return messages, nil
}

func (uc *messageUseCase) SoftDelete(id string) error {
	if err := uc.repo.SoftDelete(id); err != nil {
		return apperror.InternalServerError(err, "failed to delete message")
	}
	return nil
}

func (uc *messageUseCase) GetByConversationPaginated(convoID string, limit, page int) (*[]domain.Message, int, int, error) {
	return uc.repo.FindByConversationIDPaginated(convoID, limit, page)
}
