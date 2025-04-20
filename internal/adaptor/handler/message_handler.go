package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/internal/usecase/message"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
)

type messageHandler struct {
	msgUseCase message.MessageUseCase
	dto        dto.MessageDto
}

func NewMessageHandler(msgUseCase message.MessageUseCase, dto dto.MessageDto) *messageHandler {
	return &messageHandler{
		msgUseCase: msgUseCase,
		dto:        dto,
	}
}

// CreateMessage godoc
//
//	@summary 		CreateMessage
//	@description	Send a new message in a conversation
//	@tags 			message
//	@Security		Bearer
//	@produce		json
//	@param 			message	body	dto.CreateMessageRequest	true	"Message Data"
//	@response 		201	{object}	dto.SuccessResponse[dto.MessageResponse]	"Created"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /messages [post]
func (h *messageHandler) HandleCreateMessage(c *fiber.Ctx) error {
	body := new(dto.CreateMessageRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("failed to retrieve user from context"), "unable to retrieve user from context")
	}

	message, err := h.msgUseCase.Create(user.ID, *body)
	if err != nil {
		return err
	}
	respData, err := h.dto.ToResponse(message)
	if err != nil {
		return err
	}
	resp := dto.Success(respData)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GetMessage godoc
//
//	@summary 		GetMessage
//	@description	Get a message by ID
//	@tags 			message
//	@Security		Bearer
//	@produce		json
//	@Param 			id	path	string	true	"Message ID"
//	@response 		200	{object}	dto.SuccessResponse[dto.MessageResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		404	{object}	dto.ErrorResponse	"Not Found"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /messages/{id} [get]
func (h *messageHandler) HandleGetMessage(c *fiber.Ctx) error {
	id := c.Params("id")

	message, err := h.msgUseCase.GetByID(id)
	if err != nil {
		return err
	}

	respData, err := h.dto.ToResponse(message)
	if err != nil {
		return err
	}
	resp := dto.Success(respData)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListMessages godoc
//
//	@summary		List Messages
//	@description	List all messages by conversation ID
//	@tags			message
//	@Security		Bearer
//	@produce		json
//	@Param			conversationID	path	string	true	"Conversation ID"
//	@Param			limit	query	int	false	"Number of messages per page"
//	@Param			page	query	int	false	"Page number"
//	@response		200	{object}	dto.PaginationResponse[dto.MessageResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /conversations/{conversationID}/messages [get]
func (h *messageHandler) HandleListMessagesByConversation(c *fiber.Ctx) error {
	convoID := c.Params("conversation_id")
	page, limit := extractPaginationControl(c)

	messages, last, total, err := h.msgUseCase.GetByConversationPaginated(convoID, limit, page)
	if err != nil {
		return err
	}

	respData, err := h.dto.ToResponseList(*messages)
	if err != nil {
		return apperror.InternalServerError(err, "failed to create message response data")
	}

	resp := dto.SuccessPagination(*respData, page, last, limit, total)
	return c.Status(200).JSON(resp)
}
