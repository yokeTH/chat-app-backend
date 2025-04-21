package user

import (
	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
	"github.com/yokeTH/chat-app-backend/internal/domain"
)

type UserRepository interface {
	GetUserByID(id string) (*domain.User, error)
	GetUserByProvider(provider, providerID string) (*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
	UpdateUserInfo(userID string, updatedData dto.UpdateUserRequest) error
	SetIsOnline(userID string, isOnline bool) error
	ListUser(page, limit int) (*[]domain.User, int, int, error)
}

type UserUseCase interface {
	GoogleLogin(profile domain.Profile) (*domain.User, error)
	List(page, limit int) (*[]domain.User, int, int, error)
	Update(id string, updatedData dto.UpdateUserRequest) (*domain.User, error)
	SetUserOnline(id string) error
	SetUserOffline(id string) error
	GetGoogleProfile(googleID string) (*domain.User, error)
}
