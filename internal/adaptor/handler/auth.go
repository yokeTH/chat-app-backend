package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/chat-app-backend/internal/domain"
	"github.com/yokeTH/chat-app-backend/internal/usecase/user"
	"github.com/yokeTH/chat-app-backend/pkg/apperror"
)

type authHandler struct {
	userUseCase user.UserUseCase
}

func NewAuthHandler(userUC user.UserUseCase) *authHandler {
	return &authHandler{
		userUseCase: userUC,
	}
}

func (a *authHandler) HandleGoogleLogin(c *fiber.Ctx) error {
	profile, ok := c.Locals("profile").(domain.Profile)
	if !ok {
		return apperror.InternalServerError(errors.New("get profile error"), "get profile error")
	}
	user, err := a.userUseCase.GoogleLogin(profile)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
