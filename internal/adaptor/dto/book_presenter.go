package dto

import (
	"github.com/yokeTH/chat-app-backend/internal/domain"
)

type bookDto struct{}

type BookDto interface {
	ToResponse(book *domain.Book) BookResponse
	ToResponseList(books []domain.Book) []BookResponse
}

func NewBookDto() *bookDto {
	return &bookDto{}
}

func (p *bookDto) ToResponse(book *domain.Book) BookResponse {
	return BookResponse{
		ID:     book.ID,
		Author: book.Author,
		Title:  book.Title,
	}
}

func (p *bookDto) ToResponseList(books []domain.Book) []BookResponse {
	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = p.ToResponse(&book)
	}
	return response
}

type CreateBookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

type BookResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type UpdateBookRequest struct {
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
}
