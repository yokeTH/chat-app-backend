package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BucketType string

const (
	PublicBucketType  BucketType = "PUBLIC"
	PrivateBucketType BucketType = "PRIVATE"
)

type File struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	MessageID string    `gorm:"size:36;not null;index"`
	Key       string    `gorm:"size:255;not null"`
	MimeType  string    `gorm:"size:100"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	Message Message `gorm:"foreignKey:MessageID"`
}

func (a *File) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
