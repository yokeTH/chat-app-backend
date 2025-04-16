package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/user"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type userHandler struct {
	userUC user.UserUseCase
	dto    dto.UserDto
}

func NewUserHandler(userUC user.UserUseCase, dto dto.UserDto) *userHandler {
	return &userHandler{
		userUC: userUC,
		dto:    dto,
	}
}

func (h *userHandler) HandleListUser(c *fiber.Ctx) error {
	page, limit := extractPaginationControl(c)
	users, last, total, err := h.userUC.List(page, limit)
	if err != nil {
		return err
	}

	respData := h.dto.ToResponseList(*users)
	resp := dto.SuccessPagination(*respData, page, last, limit, total)

	return c.Status(200).JSON(resp)
}

func (h *userHandler) HandleUpdateUser(c *fiber.Ctx) error {
	body := new(dto.UpdateUserRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("failed to retrieve user from context"), "unable to retrieve user from context")
	}

	user, err := h.userUC.Update(user.ID, *body)
	if err != nil {
		return err
	}

	respData := h.dto.ToResponse(user)
	resp := dto.Success(respData)

	return c.Status(200).JSON(resp)
}
