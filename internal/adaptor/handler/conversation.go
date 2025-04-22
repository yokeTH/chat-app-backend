package handler

import (
	"errors"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/websocket"
	"github.com/yokeTH/chat-app-backend/internal/domain"
	"github.com/yokeTH/chat-app-backend/internal/usecase/conversation"
	"github.com/yokeTH/chat-app-backend/internal/usecase/message"
	"github.com/yokeTH/chat-app-backend/pkg/apperror"
)

type conversationHandler struct {
	convUC     conversation.ConversationUseCase
	msgUC      message.MessageUseCase
	dto        dto.ConversationDto
	mServer    websocket.MessageServer
	messageDto dto.MessageDto
}

func NewConversationHandler(convUC conversation.ConversationUseCase, dto dto.ConversationDto, mServer websocket.MessageServer, msgUC message.MessageUseCase, messageDto dto.MessageDto) *conversationHandler {
	return &conversationHandler{
		convUC:     convUC,
		msgUC:      msgUC,
		dto:        dto,
		mServer:    mServer,
		messageDto: messageDto,
	}
}

// GetConversations godoc
//
//	@summary		GetConversation
//	@description	list conversations
//	@tags			conversation
//	@Security		Bearer
//	@produce		json
//	@Param			limit	query	int	false	"Number of history to be retrieved"
//	@Param			page	query	int	false	"Page to retrieved"
//	@response		200	{object}	dto.PaginationResponse[dto.ConversationResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /conversations [get]
func (c *conversationHandler) HandleListConversation(ctx *fiber.Ctx) error {
	page, limit := extractPaginationControl(ctx)
	user, ok := ctx.Locals("user").(*domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("failed to retrieve user from context"), "unable to retrieve user from context")
	}

	conversations, last, total, err := c.convUC.GetUserConversations(user.ID, limit, page)
	if err != nil {
		return err
	}

	respData, err := c.dto.ToResponseList(*conversations)
	if err != nil {
		return apperror.InternalServerError(err, "failed to create response data")
	}
	resp := dto.SuccessPagination(*respData, page, last, limit, total)

	return ctx.JSON(resp)
}

// CreateNewConversation godoc
//
//	@summary		Create conversation
//	@description	create new conversation
//	@tags			conversation
//	@Security		Bearer
//	@accept			json
//	@produce 		json
//	@param			conversation	body 	dto.CreateConversationRequest	true	"conversation data"
//	@success 		201	{object}	dto.SuccessResponse[dto.ConversationResponse]	"Created"
//	@failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@failure 		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /conversations [post]
func (c *conversationHandler) HandleCreateConversation(ctx *fiber.Ctx) error {
	body := new(dto.CreateConversationRequest)
	if err := ctx.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, "invalid body")
	}
	user, ok := ctx.Locals("user").(*domain.User)
	if !ok {
		return apperror.InternalServerError(errors.New("get profile error"), "get profile error")
	}
	conversation, err := c.convUC.CreateConversation(body.Members, user.ID, body.Name)
	if err != nil {
		return err
	}

	respData, err := c.dto.ToResponse(conversation)
	if err != nil {
		return apperror.InternalServerError(err, "failed to create response data")
	}
	resp := dto.Success(respData)

	system, err := c.msgUC.CreateSystemMessage(conversation.ID, "chat has been created")
	if err != nil {
		return err
	}

	createdMessageResponse, err := c.messageDto.ToResponse(system)
	if err != nil {
		log.Printf("failed to transform to dto: %v", err)
		return err
	}

	payloadResponse, err := json.Marshal(createdMessageResponse)
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return err
	}

	createdMessageJson, err := json.Marshal(websocket.WebSocketMessage{
		Event:     websocket.EventTypeMessage,
		Payload:   payloadResponse,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		return apperror.InternalServerError(err, err.Error())
	}

	if err := c.mServer.BroadcastToMembersInConversation(conversation.ID, createdMessageJson); err != nil {
		return apperror.InternalServerError(err, "broadcast error")
	}

	return ctx.Status(201).JSON(resp)
}

// GetConversations godoc
//
//	@summary		Get Conversation by id
//	@description	Get Conversation by id
//	@tags			conversation
//	@Security		Bearer
//	@produce		json
//	@Param			limit	query	int		false	"Number of history to be retrieved"
//	@Param			page	query	int		false	"Page to retrieved"
//	@Param			id		path	string	true	"conversation id"
//	@response		200	{object}	dto.SuccessResponse[dto.ConversationResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /conversations/{id} [get]
func (c *conversationHandler) HandleGetConversation(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	conversation, err := c.convUC.GetConversation(id)
	if err != nil {
		return err
	}

	respData, err := c.dto.ToResponse(conversation)
	if err != nil {
		return apperror.InternalServerError(err, "failed to create response data")
	}
	resp := dto.Success(*respData)

	return ctx.JSON(resp)
}
