package repository

import (
	"errors"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/db"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create message")
	}
	return nil
}

func (r *messageRepository) FindByID(id string) (*domain.Message, error) {
	var message domain.Message

	if err := r.db.Preload("Sender").Preload("Conversation").Preload("Attachments").Preload("Reactions").First(&message, "id = ?", id).Error; err != nil {
		return nil, apperror.InternalServerError(err, "failed to find message by id")
	}

	return &message, nil
}

func (r *messageRepository) FindByConversationID(conversationID string) (*[]domain.Message, error) {
	var messages []domain.Message

	if err := r.db.
		Preload("Sender").
		Preload("Attachments").
		Preload("Reactions").
		Where("conversation_id = ? AND is_deleted = false", conversationID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, apperror.InternalServerError(err, "failed to find by conversation id")
	}
	return &messages, nil
}

func (r *messageRepository) Update(message *domain.Message) error {
	if err := r.db.Save(message).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update message")
	}
	return nil
}

func (r *messageRepository) Delete(id string) error {
	if err := r.db.Delete(&domain.Message{}, "id = ?", id).Error; err != nil {
		return apperror.InternalServerError(err, "failed to delete message")
	}
	return nil
}

func (r *messageRepository) SoftDelete(id string) error {
	if err := r.db.Model(&domain.Message{}).Where("id = ?", id).Update("is_deleted", true).Error; err != nil {
		return apperror.InternalServerError(err, "failed to delete message")
	}
	return nil
}

func (r *messageRepository) FindByConversationIDPaginated(convoID string, limit, page int) (*[]domain.Message, int, int, error) {
	var messages []domain.Message
	var total, last int

	query := r.db.
		Where("conversation_id = ? AND is_deleted = false", convoID).
		Order("created_at ASC").
		Preload("Sender").
		Preload("Attachments").
		Preload("Reactions").
		Scopes(db.Paginate(&domain.Message{}, &limit, &page, &total, &last))

	if err := query.Find(&messages).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "messages not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to fetch messages")
	}

	return &messages, last, total, nil
}
