package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText   MessageType = "TEXT"
	MessageTypeSystem MessageType = "SYSTEM"
)

type Message struct {
	ID             string      `gorm:"primaryKey;type:varchar(36)"`
	ConversationID string      `gorm:"size:36;not null;index"`
	SenderID       string      `gorm:"size:36;index;default:null"`
	Content        string      `gorm:"type:text"`
	CreatedAt      time.Time   `gorm:"autoCreateTime;index"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime"`
	IsDeleted      bool        `gorm:"default:false"`
	MessageType    MessageType `gorm:"default:TEXT"`

	// Relationships
	Conversation Conversation `gorm:"foreignKey:ConversationID"`
	Sender       User         `gorm:"foreignKey:SenderID; default:null"`
	Attachments  []File       `gorm:"foreignKey:MessageID"`
	Reactions    []Reaction   `gorm:"foreignKey:MessageID"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	var count int64
	if err := tx.Table("conversation_members").
		Where("conversation_id = ? AND user_id = ?", m.ConversationID, m.SenderID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user is not a member of the conversation")
	}
	return nil
}
