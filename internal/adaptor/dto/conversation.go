package dto

import (
	"github.com/yokeTH/gofiber-template/internal/domain"
)

type conversationDto struct {
	userDto    UserDto
	messageDto MessageDto
}

type ConversationDto interface {
	ToResponse(conversation *domain.Conversation) (*ConversationResponse, error)
	ToResponseList(conversations []domain.Conversation) (*[]ConversationResponse, error)
}

func NewConversationDto(userDto UserDto, messageDto MessageDto) *conversationDto {
	return &conversationDto{
		userDto:    userDto,
		messageDto: messageDto,
	}
}

func (c *conversationDto) ToResponse(conversation *domain.Conversation) (*ConversationResponse, error) {
	messages, err := c.messageDto.ToResponseList(conversation.Messages)
	if err != nil {
		return nil, err
	}
	return &ConversationResponse{
		ID:       conversation.ID,
		Name:     conversation.Name,
		Members:  *c.userDto.ToResponseList(conversation.Members),
		Messages: *messages,
	}, nil
}

func (c *conversationDto) ToResponseList(conversations []domain.Conversation) (*[]ConversationResponse, error) {
	response := make([]ConversationResponse, len(conversations))
	for i, conversation := range conversations {
		resp, err := c.ToResponse(&conversation)
		if err != nil {
			return nil, err
		}
		response[i] = *resp
	}

	return &response, nil
}

type ConversationResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Members  []UserResponse    `json:"members"`
	Messages []MessageResponse `json:"messages"`
}
