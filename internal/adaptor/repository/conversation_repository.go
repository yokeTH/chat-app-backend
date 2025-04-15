package repository

import (
	"errors"
	"fmt"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/db"
	"gorm.io/gorm"
)

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *conversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) GetUserConversations(userID string, limit, page int) (*[]domain.Conversation, int, int, error) {
	var conversations []domain.Conversation
	var total, last int

	if err := r.db.
		Joins("JOIN conversation_members ON conversation_members.conversation_id = conversations.id").
		Where("conversation_members.user_id = ?", userID).
		Scopes(db.Paginate(&domain.Conversation{}, &limit, &page, &total, &last)).
		Preload("Members").
		Preload("Messages").
		Find(&conversations).
		Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "fail to retrieve conversation")
	}

	return &conversations, last, total, nil
}

func (r *conversationRepository) CreateConversation(usersID []string, createdByID string, name string) (*domain.Conversation, error) {
	if len(usersID) < 2 {
		return nil, errors.New("invalid usersID cannot be empty")
	}

	isGroup := len(usersID) > 2

	conversation := &domain.Conversation{
		Name:      name,
		IsGroup:   isGroup,
		CreatedBy: createdByID,
		Members:   []domain.User{},
	}

	var users []domain.User
	if err := r.db.Where("id IN ?", usersID).Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) != len(usersID) {
		return nil, fmt.Errorf("some user IDs are invalid")
	}

	conversation.Members = users

	if err := r.db.Create(&conversation).Error; err != nil {
		return nil, err
	}

	return conversation, nil
}
