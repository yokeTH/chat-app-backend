package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationMember struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)"`
	ConversationID string    `gorm:"size:36;not null;index:idx_conversation_member,unique:idx_conversation_member"`
	UserID         string    `gorm:"size:36;not null;index:idx_conversation_member,unique:idx_conversation_member"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	IsAdmin        bool      `gorm:"default:false"`

	// Relationships
	Conversation Conversation `gorm:"foreignKey:ConversationID"`
	User         User         `gorm:"foreignKey:UserID"`
}

func (cm *ConversationMember) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == "" {
		cm.ID = uuid.New().String()
	}
	return nil
}
