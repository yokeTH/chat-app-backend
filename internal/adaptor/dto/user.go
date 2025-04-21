package dto

import (
	"time"

	"github.com/yokeTH/chat-app-backend/internal/domain"
)

type UserDto interface {
	ToResponse(user *domain.User) *UserResponse
	ToResponseList(users []domain.User) *[]UserResponse
}

type userDto struct{}

func NewUserDto() *userDto {
	return &userDto{}
}

func (u *userDto) ToResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		IsOnline:  user.IsOnline,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userDto) ToResponseList(users []domain.User) *[]UserResponse {
	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = *u.ToResponse(&user)
	}
	return &response
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar"`
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
