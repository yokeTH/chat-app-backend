package user

import (
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type UserRepository interface {
	GetUserByID(id string) (*domain.User, *apperror.AppError)
	GetUserByProvider(provider, providerID string) (*domain.User, *apperror.AppError)
	CreateUser(user *domain.User) (*domain.User, *apperror.AppError)
	UpdateUserInfo(userID string, updatedData dto.UpdateUserRequest) *apperror.AppError
	SetIsOnline(userID string, isOnline bool) *apperror.AppError
}

type UserUseCase interface {
	GoogleLogin(profile domain.Profile) (*domain.User, *apperror.AppError)
}
