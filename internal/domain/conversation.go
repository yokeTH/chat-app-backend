package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Name      string    `gorm:"size:100"`
	IsGroup   bool      `gorm:"default:false"`
	CreatedBy string    `gorm:"size:36;not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Members  []User    `gorm:"many2many:conversation_members;"`
	Messages []Message `gorm:"foreignKey:ConversationID"`
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
