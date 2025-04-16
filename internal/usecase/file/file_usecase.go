package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

type fileUseCase struct {
	fileRepo   FileRepository
	pubStorage storage.Storage
}

func NewFileUseCase(fileRepo FileRepository, pub storage.Storage) *fileUseCase {
	return &fileUseCase{
		fileRepo:   fileRepo,
		pubStorage: pub,
	}
}

func (u *fileUseCase) List(limit, page int) ([]domain.File, int, int, error) {
	return u.fileRepo.List(limit, page)
}

func (u *fileUseCase) GetByID(id int) (*domain.File, error) {
	return u.fileRepo.GetByID(id)
}

func (u *fileUseCase) CreateFile(ctx context.Context, file *multipart.FileHeader) (*domain.File, error) {
	fileData, err := file.Open()
	if err != nil {
		return nil, apperror.InternalServerError(err, "error opening file")
	}
	defer fileData.Close()

	filename := strings.ReplaceAll(file.Filename, " ", "-")
	contentType := file.Header.Get("Content-Type")
	fileKey := fmt.Sprintf("upload/%s", filename)

	fileInfo := &domain.File{
		Key: fileKey,
	}

	if err = u.pubStorage.UploadFile(ctx, fileKey, contentType, fileData); err != nil {
		return nil, apperror.InternalServerError(err, "error uploading file")
	}

	if err = u.fileRepo.Create(fileInfo); err != nil {
		return nil, apperror.InternalServerError(err, "error create file data")
	}

	return fileInfo, nil
}
