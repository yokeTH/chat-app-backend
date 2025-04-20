package conversation

import "github.com/yokeTH/gofiber-template/internal/domain"

type conversationUseCase struct {
	convRepo ConversationRepository
}

func NewConversationUseCase(convRepo ConversationRepository) *conversationUseCase {
	return &conversationUseCase{
		convRepo: convRepo,
	}
}

func (c *conversationUseCase) GetUserConversations(userID string, limit, page int) (*[]domain.Conversation, int, int, error) {
	return c.convRepo.GetUserConversations(userID, limit, page)
}

func (c *conversationUseCase) CreateConversation(usersID []string, createdByID string, name string) (*domain.Conversation, error) {
	return c.convRepo.CreateConversation(usersID, createdByID, name)
}

func (c *conversationUseCase) GetMembers(id string) (*[]domain.User, error) {
	return c.convRepo.GetMembers(id)
}
