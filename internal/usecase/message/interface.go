package message

import (
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
)

type MessageRepository interface {
	Create(message *domain.Message) error
	FindByID(id string) (*domain.Message, error)
	FindByConversationID(conversationID string) (*[]domain.Message, error)
	FindByConversationIDPaginated(convoID string, limit, page int) (*[]domain.Message, int, int, error)
	Update(message *domain.Message) error
	Delete(id string) error
	SoftDelete(id string) error
}

type MessageUseCase interface {
	Create(senderID string, req dto.CreateMessageRequest) (*domain.Message, error)
	GetByID(id string) (*domain.Message, error)
	GetByConversationID(convoID string) (*[]domain.Message, error)
	GetByConversationPaginated(convoID string, limit, page int) (*[]domain.Message, int, int, error)
	SoftDelete(id string) error
}
