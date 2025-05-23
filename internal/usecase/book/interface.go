package book

import (
	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
	"github.com/yokeTH/chat-app-backend/internal/domain"
)

type BookRepository interface {
	Create(book *domain.Book) error
	GetByID(id int) (*domain.Book, error)
	List(limit, page int) ([]domain.Book, int, int, error)
	Update(id int, book *dto.UpdateBookRequest) (*domain.Book, error)
	Delete(id int) error
}

type BookUseCase interface {
	Create(book *domain.Book) error
	GetByID(id int) (*domain.Book, error)
	List(limit, page int) ([]domain.Book, int, int, error)
	Update(id int, book *dto.UpdateBookRequest) (*domain.Book, error)
	Delete(id int) error
}
