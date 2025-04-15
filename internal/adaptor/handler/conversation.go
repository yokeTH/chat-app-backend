package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/conversation"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type conversationHandler struct {
	convUC conversation.ConversationUseCase
}

func NewConversationHandler(convUC conversation.ConversationUseCase) *conversationHandler {
	return &conversationHandler{
		convUC: convUC,
	}
}

func (c *conversationHandler) HandleListConversation(ctx *fiber.Ctx) error {
	page, limit := extractPaginationControl(ctx)
	user, ok := ctx.Locals("user").(domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("get user error"), "get error error")
	}

	conversations, last, total, err := c.convUC.GetUserConversations(user.ID, limit, page)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.SuccessPagination(*conversations, page, last, limit, total))
}
