package conversation

import "github.com/yokeTH/chat-app-backend/internal/domain"

type ConversationRepository interface {
	GetUserConversations(userID string, limit, page int) (*[]domain.Conversation, int, int, error)
	CreateConversation(usersID []string, createdByID string, name string) (*domain.Conversation, error)
	GetMembers(id string) (*[]domain.User, error)
	GetConversation(id string) (*domain.Conversation, error)
	GetUserNotInConversations(userID string, limit, page int) (*[]domain.Conversation, int, int, error)
	AddMemberToConversation(conversationID, userID string) error
}

type ConversationUseCase interface {
	GetUserConversations(userID string, limit, page int) (*[]domain.Conversation, int, int, error)
	CreateConversation(usersID []string, createdByID string, name string) (*domain.Conversation, error)
	GetMembers(id string) (*[]domain.User, error)
	GetConversation(id string) (*domain.Conversation, error)
	AddMember(conversationID, userID string) error
}
