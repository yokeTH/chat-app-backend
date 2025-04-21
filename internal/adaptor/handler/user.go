package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
	"github.com/yokeTH/chat-app-backend/internal/domain"
	"github.com/yokeTH/chat-app-backend/internal/usecase/user"
	"github.com/yokeTH/chat-app-backend/pkg/apperror"
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

// GetUsers godoc
//
//	@summary		GetUsers
//	@description	get users
//	@tags			user
//	@Security		Bearer
//	@produce		json
//	@Param			limit	query	int	false	"Number of history to be retrieved"
//	@Param			page	query	int	false	"Page to retrieved"
//	@response		200	{object}	dto.PaginationResponse[dto.UserResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /users [get]
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

// UpdateUser godoc
//
//	@summary 		UpdateUser
//	@description	update user data
//	@tags 			user
//	@Security		Bearer
//	@produce		json
//	@Param 			id		path	string					true	"User ID"
//	@param 			user	body	dto.UpdateUserRequest	true	"User Data"
//	@response 		200	{object}	dto.SuccessResponse[dto.UserResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@response		403	{object}	dto.ErrorResponse	"Forbidden"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /users/{id} [patch]
func (h *userHandler) HandleUpdateUser(c *fiber.Ctx) error {
	body := new(dto.UpdateUserRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("failed to retrieve user from context"), "unable to retrieve user from context")
	}

	id := c.Params("id")

	// perm check
	if user.ID != id {
		return apperror.ForbiddenError(fmt.Errorf("no permission to edit user information: context user id = %s, params id = %s", user.ID, id), "no permission to edit user information")
	}

	user, err := h.userUC.Update(id, *body)
	if err != nil {
		return err
	}

	respData := h.dto.ToResponse(user)
	resp := dto.Success(respData)

	return c.Status(200).JSON(resp)
}

// GetMyUser godoc
//
//	@summary 		Get My User
//	@description	get my user data
//	@tags 			user
//	@Security		Bearer
//	@produce		json
//	@response 		200	{object}	dto.SuccessResponse[dto.UserResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /users/me [get]
func (h *userHandler) HandleGetMe(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("failed to retrieve user from context"), "unable to retrieve user from context")
	}

	respData := h.dto.ToResponse(user)
	resp := dto.Success(respData)

	return c.Status(200).JSON(resp)
}
