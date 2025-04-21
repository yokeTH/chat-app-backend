package repository

import (
	"errors"

	"github.com/yokeTH/chat-app-backend/internal/adaptor/dto"
	"github.com/yokeTH/chat-app-backend/internal/domain"
	"github.com/yokeTH/chat-app-backend/pkg/apperror"
	"github.com/yokeTH/chat-app-backend/pkg/db"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) Create(book *domain.Book) error {
	if err := r.db.Create(book).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create book")
	}
	return nil
}

func (r *bookRepository) GetByID(id int) (*domain.Book, error) {
	book := &domain.Book{}
	if err := r.db.First(book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "book not found")
		}
		return nil, apperror.InternalServerError(err, "failed to get book")
	}
	return book, nil
}

func (r *bookRepository) List(limit, page int) ([]domain.Book, int, int, error) {
	var books []domain.Book
	var total, last int

	if err := r.db.Scopes(db.Paginate(domain.Book{}, &limit, &page, &total, &last)).Find(&books).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "books not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to get books")
	}
	return books, last, total, nil
}

func (r *bookRepository) Update(id int, updateRequest *dto.UpdateBookRequest) (*domain.Book, error) {
	var book domain.Book
	if err := r.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "book not found")
		}
		return nil, err
	}

	if err := r.db.Model(&book).Updates(updateRequest).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *bookRepository) Delete(id int) error {
	if err := r.db.Delete(&domain.Book{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError(err, "book not found")
		}
		return apperror.InternalServerError(err, "failed to delete book")
	}
	return nil
}
