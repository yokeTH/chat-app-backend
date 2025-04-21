package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
	"github.com/yokeTH/chat-app-backend/internal/usecase/file"
	"github.com/yokeTH/chat-app-backend/pkg/apperror"
)

type fileHandler struct {
	fileUseCase file.FileUseCase
	dto         dto.FileDto
}

func NewFileHandler(uc file.FileUseCase, dto dto.FileDto) *fileHandler {
	return &fileHandler{
		fileUseCase: uc,
		dto:         dto,
	}
}

// CreatePublicFile godoc
//
//	@summary		CreatePublicFile
//	@description	create public file by upload file multipart-form field name file
//	@tags			file
//	@accept			x-www-form-urlencoded
//	@produce 		json
//	@param			file	formData 	file	true "file data"
//	@success 		201	{object}	dto.SuccessResponse[dto.FileResponse]	"Created"
//	@failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@failure 		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /files [post]
func (h *fileHandler) CreateFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequestError(err, "invalid file")
	}

	fileData, err := h.fileUseCase.CreateFile(c.Context(), file)
	if err != nil {
		return err
	}

	response, err := h.dto.ToResponse(*fileData)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(dto.Success(response))
}

// GetFiles godoc
//
//	@summary		GetBooks
//	@description	get files information
//	@tags			file
//	@produce		json
//	@Param			limit	query	int	false	"Number of history to be retrieved"
//	@Param			page	query	int	false	"Page to retrieved"
//	@response		200	{object}	dto.PaginationResponse[dto.FileResponse]	"OK"
//	@response		400	{object}	dto.ErrorResponse	"Bad Request"
//	@response		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /files [get]
func (h *fileHandler) List(c *fiber.Ctx) error {
	page, limit := extractPaginationControl(c)
	files, last, total, err := h.fileUseCase.List(limit, page)
	if err != nil {
		return err
	}

	response, err := h.dto.ToResponseList(files)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(dto.SuccessPagination(*response, page, last, limit, total))
}

// GetFile godoc
//
//	@summary		Get file url
//	@description	get file information and url
//	@tags			file
//	@produce		json
//	@Param			id	path	int	true	"file id"
//	@success 		200	{object}	dto.SuccessResponse[dto.FileResponse]	"OK"
//	@failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@failure 		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router /files/{id} [get]
func (h *fileHandler) GetInfo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "invalid id")
	}

	file, err := h.fileUseCase.GetByID(id)
	if err != nil {
		return err
	}

	response, err := h.dto.ToResponse(*file)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(dto.Success(response))
}
