package repository

import (
	"fmt"
	"sort"

	"github.com/yokeTH/chat-app-backend/internal/domain"
	"github.com/yokeTH/chat-app-backend/pkg/apperror"
	"github.com/yokeTH/chat-app-backend/pkg/db"
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
		Find(&conversations).
		Error; err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "fail to retrieve conversation")
	}

	for i := range conversations {
		var lastMessage []domain.Message
		if err := r.db.
			Where("conversation_id = ?", conversations[i].ID).
			Order("created_at DESC").
			// Limit(10).
			Preload("Sender").
			Find(&lastMessage).Error; err != nil {
			return nil, 0, 0, apperror.InternalServerError(err, "fail to retrieve last message")
		}
		sort.Slice(lastMessage, func(i, j int) bool {
			return lastMessage[i].CreatedAt.Before(lastMessage[j].CreatedAt)
		})
		conversations[i].Messages = lastMessage
	}

	return &conversations, last, total, nil
}

func (r *conversationRepository) CreateConversation(usersID []string, createdByID string, name string) (*domain.Conversation, error) {
	if len(usersID) < 2 {
		return nil, apperror.BadRequestError(fmt.Errorf("validate create conversation failed"), "users id must be more than 1")
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
	fmt.Println("FIND ID IN:", users)
	if len(users) != len(usersID) {
		return nil, apperror.BadRequestError(fmt.Errorf("some user IDs are invalid in create conversation repository"), "some user IDs are invalid")
	}

	conversation.Members = users

	if err := r.db.Create(&conversation).Error; err != nil {
		return nil, err
	}

	return conversation, nil
}

func (r *conversationRepository) GetMembers(id string) (*[]domain.User, error) {
	var conversation domain.Conversation
	if err := r.db.
		Where("id = ?", id).
		Preload("Members").
		First(&conversation).
		Error; err != nil {
		return nil, apperror.InternalServerError(err, "failed to retrieve conversation")
	}
	return &conversation.Members, nil
}

func (r *conversationRepository) GetConversation(id string) (*domain.Conversation, error) {
	var conversation domain.Conversation

	if err := r.db.
		Where("id = ?", id).
		Preload("Members").
		First(&conversation).
		Error; err != nil {
		return nil, apperror.InternalServerError(err, "fail to retrieve conversation")
	}

	var lastMessage []domain.Message
	if err := r.db.
		Where("conversation_id = ?", conversation.ID).
		Order("created_at DESC").
		// Limit(10).
		Preload("Sender").
		Find(&lastMessage).Error; err != nil {
		return nil, apperror.InternalServerError(err, "fail to retrieve last message")
	}
	sort.Slice(lastMessage, func(i, j int) bool {
		return lastMessage[i].CreatedAt.Before(lastMessage[j].CreatedAt)
	})
	conversation.Messages = lastMessage

	return &conversation, nil
}
