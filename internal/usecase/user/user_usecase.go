package user

import (
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
