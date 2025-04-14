package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reaction struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	MessageID string    `gorm:"size:36;not null;index:idx_reaction_unique,unique:idx_reaction_unique"`
	UserID    string    `gorm:"size:36;not null;index:idx_reaction_unique,unique:idx_reaction_unique"`
	Emoji     string    `gorm:"size:10;not null;index:idx_reaction_unique,unique:idx_reaction_unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	Message Message `gorm:"foreignKey:MessageID"`
	User    User    `gorm:"foreignKey:UserID"`
}

func (r *Reaction) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
