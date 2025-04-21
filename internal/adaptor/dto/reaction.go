package dto

import "github.com/yokeTH/chat-app-backend/internal/domain"

type reactionDto struct {
	userDto UserDto
}

type ReactionDto interface {
	ToResponse(e *domain.Reaction) *ReactionResponse
	ToResponseList(es []domain.Reaction) *[]ReactionResponse
}

func NewReactionDto(userDto UserDto) *reactionDto {
	return &reactionDto{
		userDto: userDto,
	}
}

func (r *reactionDto) ToResponse(e *domain.Reaction) *ReactionResponse {
	return &ReactionResponse{
		Emoji: e.Emoji,
		User:  *r.userDto.ToResponse(&e.User),
	}
}

func (r *reactionDto) ToResponseList(es []domain.Reaction) *[]ReactionResponse {
	response := make([]ReactionResponse, len(es))
	for i, e := range es {
		response[i] = *r.ToResponse(&e)
	}
	return &response
}

type ReactionResponse struct {
	Emoji string       `json:"emoji"`
	User  UserResponse `json:"user"`
}
