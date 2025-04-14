package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	MessageID string    `gorm:"size:36;not null;index"`
	URL       string    `gorm:"size:255;not null"`
	Size      int64     `json:"size"`
	MimeType  string    `gorm:"size:100"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	Message Message `gorm:"foreignKey:MessageID"`
}

func (a *Attachment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
