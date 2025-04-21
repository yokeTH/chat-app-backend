package dto

import (
	"time"

	"github.com/yokeTH/chat-app-backend/internal/domain"
)

type MessageDto interface {
	ToResponse(e *domain.Message) (*MessageResponse, error)
	ToResponseList(es []domain.Message) (*[]MessageResponse, error)
}

type messageDto struct {
	fileDto     FileDto
	reactionDto ReactionDto
	userDto     UserDto
}

func NewMessageDto(fileDto FileDto, reactionDto ReactionDto, userDto UserDto) *messageDto {
	return &messageDto{
		fileDto:     fileDto,
		reactionDto: reactionDto,
		userDto:     userDto,
	}
}

func (m *messageDto) ToResponse(e *domain.Message) (*MessageResponse, error) {
	attachments, err := m.fileDto.ToResponseList(e.Attachments)
	if err != nil {
		return nil, err
	}

	return &MessageResponse{
		ID:             e.ID,
		Content:        e.Content,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		Sender:         *m.userDto.ToResponse(&e.Sender),
		ConversationID: e.ConversationID,
		Attachments:    *attachments,
		Reactions:      *m.reactionDto.ToResponseList(e.Reactions),
	}, nil
}

func (m *messageDto) ToResponseList(es []domain.Message) (*[]MessageResponse, error) {
	response := make([]MessageResponse, len(es))
	for i, e := range es {
		resp, err := m.ToResponse(&e)
		if err != nil {
			return nil, err
		}
		response[i] = *resp
	}
	return &response, nil
}

type MessageResponse struct {
	ID             string       `json:"id"`
	Content        string       `json:"content"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Sender         UserResponse `json:"sender"`
	ConversationID string       `json:"conversation_id"`
	// Sender      UserResponse       `json:"senderId,omitempty"`
	Attachments []FileResponse     `json:"attachments"`
	Reactions   []ReactionResponse `json:"reactions"`
}

type CreateMessageRequest struct {
	ConversationID string `json:"conversation_id" validate:"required,uuid4"`
	Content        string `json:"content"`
}
