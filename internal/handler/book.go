package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/core/domain"
	"github.com/yokeTH/gofiber-template/internal/core/port"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/response"
)

type BookHandler struct {
	BookService port.BookService
}

func NewBookHandler(bookService port.BookService) port.BookHandler {
	return &BookHandler{
		BookService: bookService,
	}
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	body := new(domain.Book)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	if err := h.BookService.CreateBook(body); err != nil {
		return apperror.InternalServerError(err, err.Error())
	}

	return c.JSON(response.SuccessResponse[domain.Book]{Data: *body})
}

func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book, err := h.BookService.GetBook(id)
	if err != nil {
		return apperror.InternalServerError(err, err.Error())
	}

	return c.JSON(response.SuccessResponse[domain.Book]{Data: *book})
}

func (h *BookHandler) GetBooks(c *fiber.Ctx) error {
	books, err := h.BookService.GetBooks()
	if err != nil {
		return apperror.InternalServerError(err, err.Error())
	}

	convertedBooks := make([]domain.Book, len(books))
	for i, book := range books {
		convertedBooks[i] = *book
	}
	return c.JSON(response.SuccessResponse[[]domain.Book]{Data: convertedBooks})
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	body := new(domain.Book)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book, err := h.BookService.UpdateBook(id, body)
	if err != nil {
		return apperror.InternalServerError(err, err.Error())
	}

	return c.JSON(response.SuccessResponse[domain.Book]{Data: *book})
}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	if err := h.BookService.DeleteBook(id); err != nil {
		return apperror.InternalServerError(err, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
