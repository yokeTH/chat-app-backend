package message

import (
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type messageUseCase struct {
	repo MessageRepository
}

func NewMessageUseCase(repo MessageRepository) *messageUseCase {
	return &messageUseCase{
		repo: repo,
	}
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

	message, err := uc.repo.FindByID(message.ID)
	if err != nil {
		return nil, apperror.InternalServerError(err, "failed to get message")
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
