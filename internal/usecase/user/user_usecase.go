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

func (u *userUseCase) GoogleLogin(profile domain.Profile) (*domain.User, *apperror.AppError) {

	user := domain.User{
		Name:       profile.Name,
		Email:      profile.Email,
		AvatarURL:  profile.Picture,
		Provider:   "GOOGLE",
		ProviderID: profile.Sub,
	}

	createdUser, createdErr := u.userRepo.CreateUser(&user)
	foundUser, foundErr := u.userRepo.GetUserByProvider("GOOGLE", profile.Sub)

	if foundErr != nil && createdErr != nil {
		return nil, apperror.BadRequestError(foundErr, "login error")
	} else if foundErr == nil {
		return foundUser, nil
	} else {
		return createdUser, nil
	}

}
