package dto

import (
	"time"

	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

type fileDto struct {
	public storage.Storage
}

type FileDto interface {
	ToResponse(domain.File) (*FileResponse, error)
	ToResponseList(files []domain.File) (*[]FileResponse, error)
}

func NewFileDto(pub storage.Storage) *fileDto {
	return &fileDto{
		public: pub,
	}
}

func (f *fileDto) ToResponse(file domain.File) (*FileResponse, error) {
	url, err := f.public.GetPublicUrl(file.Key)
	if err != nil {
		return nil, err
	}

	return &FileResponse{
		ID:        file.ID,
		Url:       url,
		CreatedAt: &file.CreatedAt,
		MimeType:  file.MimeType,
	}, nil

}

func (f *fileDto) ToResponseList(files []domain.File) (*[]FileResponse, error) {
	response := make([]FileResponse, len(files))
	for i, file := range files {
		resp, err := f.ToResponse(file)
		if err != nil {
			return nil, err
		}
		response[i] = *resp
	}
	return &response, nil
}

type FileResponse struct {
	ID        string     `json:"id"`
	Url       string     `json:"url,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	MimeType  string     `json:"mime_type,omitempty"`
}
