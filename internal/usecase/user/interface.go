package user

import (
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
)

type UserRepository interface {
	GetUserByID(id string) (*domain.User, error)
	GetUserByProvider(provider, providerID string) (*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
	UpdateUserInfo(userID string, updatedData dto.UpdateUserRequest) error
	SetIsOnline(userID string, isOnline bool) error
}

type UserUseCase interface {
	GoogleLogin(profile domain.Profile) (*domain.User, error)
}
