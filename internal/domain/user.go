package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash string    `gorm:"size:255;not null"`
	AvatarURL    string    `gorm:"size:255"`
	IsOnline     bool      `gorm:"default:false;index"` // to be use redis
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	Provider   string `gorm:"size:50;uniqueIndex:composite_provider"`
	ProviderID string `gorm:"size:100;uniqueIndex:composite_provider"`

	// Relationships
	Conversations []Conversation `gorm:"many2many:conversation_members;"`
	Messages      []Message      `gorm:"foreignKey:SenderID"`
	Reactions     []Reaction     `gorm:"foreignKey:UserID"`
}

type Profile struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
