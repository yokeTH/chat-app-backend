package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)"`
	ConversationID string    `gorm:"size:36;not null;index"`
	SenderID       string    `gorm:"size:36;not null;index"`
	Content        string    `gorm:"type:text"`
	SentAt         time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	IsDeleted      bool      `gorm:"default:false"`

	// Relationships
	Conversation Conversation `gorm:"foreignKey:ConversationID"`
	Sender       User         `gorm:"foreignKey:SenderID"`
	Attachments  []Attachment `gorm:"foreignKey:MessageID"`
	Reactions    []Reaction   `gorm:"foreignKey:MessageID"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
