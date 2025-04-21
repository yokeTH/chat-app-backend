package file

import (
	"context"
	"mime/multipart"

	"github.com/yokeTH/chat-app-backend/internal/domain"
)

type FileRepository interface {
	Create(file *domain.File) error
	List(limit, page int) ([]domain.File, int, int, error)
	GetByID(id int) (*domain.File, error)
}

type FileUseCase interface {
	CreateFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error)
	List(limit, page int) ([]domain.File, int, int, error)
	GetByID(id int) (*domain.File, error)
}
