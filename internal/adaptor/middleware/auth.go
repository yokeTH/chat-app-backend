package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/user"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type authMiddleware struct {
	userUseCase user.UserUseCase
}

func NewAuthMiddleware(userUseCase user.UserUseCase) *authMiddleware {
	return &authMiddleware{userUseCase: userUseCase}
}

func (a *authMiddleware) Auth(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return apperror.UnauthorizedError(errors.New("request without authorization header"), "Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return apperror.UnauthorizedError(errors.New("invalid authorization header"), "Authorization header is invalid")
	}

	token := authHeader[7:]

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return apperror.UnauthorizedError(err, "failed to create request to Google OAuth")
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return apperror.UnauthorizedError(err, "failed to get profile from Google OAuth")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		return apperror.UnauthorizedError(fmt.Errorf("non-200 response from Google OAuth: %s", resp.Status), "Failed to get profile from Google OAuth")
	}

	var profile domain.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return apperror.UnauthorizedError(err, "failed to decode profile from Google OAuth")
	}

	ctx.Locals("profile", profile)

	user, err := a.userUseCase.GoogleLogin(profile)
	if err != nil {
		return err
	}
	ctx.Locals("user", user)
	return ctx.Next()
}
