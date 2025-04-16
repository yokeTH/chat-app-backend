package repository

import (
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByID(id string) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.NotFoundError(err, "user not found")
		}
		return nil, apperror.InternalServerError(err, "failed to find user")
	}
	return &user, nil
}

func (r *userRepository) GetUserByProvider(provider, providerID string) (*domain.User, error) {
	var user domain.User

	if err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.NotFoundError(err, "user not found")
		}
		return nil, apperror.InternalServerError(err, "failed to find user")
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *domain.User) (*domain.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, apperror.InternalServerError(err, "failed to create user")
	}
	return user, nil
}

func (r *userRepository) UpdateUserInfo(userID string, updatedData dto.UpdateUserRequest) error {
	if err := r.db.
		Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(updatedData).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update user info")
	}
	return nil
}

func (r *userRepository) SetIsOnline(userID string, isOnline bool) error {
	user := domain.User{
		IsOnline: isOnline,
	}

	if err := r.db.
		Model(&domain.User{}).
		Where("id = ?", userID).
		Select("IsOnline").
		Updates(&user).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update user info")
	}

	return nil
}
