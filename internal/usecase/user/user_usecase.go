package user

import (
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type userUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *userUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) GoogleLogin(profile domain.Profile) (*domain.User, error) {
	user, err := u.userRepo.GetUserByProvider("GOOGLE", profile.Sub)
	if err == nil {
		return user, nil
	}

	newUser := domain.User{
		Name:       profile.Name,
		Email:      profile.Email,
		AvatarURL:  profile.Picture,
		Provider:   "GOOGLE",
		ProviderID: profile.Sub,
	}

	createdUser, err := u.userRepo.CreateUser(&newUser)
	if err != nil {
		return nil, apperror.BadRequestError(err, "failed to create user")
	}

	return createdUser, nil
}

func (u *userUseCase) List(page, limit int) (*[]domain.User, int, int, error) {
	return u.userRepo.ListUser(page, limit)
}

func (u *userUseCase) Update(id string, updatedData dto.UpdateUserRequest) (*domain.User, error) {
	if err := u.userRepo.UpdateUserInfo(id, updatedData); err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) SetUserOnline(id string) error {
	return u.userRepo.SetIsOnline(id, true)
}

func (u *userUseCase) SetUserOffline(id string) error {
	return u.userRepo.SetIsOnline(id, false)
}

func (u *userUseCase) GetGoogleProfile(googleID string) (*domain.User, error) {
	return u.userRepo.GetUserByProvider("GOOGLE", googleID)
}
