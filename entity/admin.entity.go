package entity

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	Id             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Username       string     `gorm:"unique" json:"username"`
	DisplayName    string     `json:"display_name"`
	ProfilePicture string     `json:"profile_picture"`
	Email          string     `json:"email"`
	Password       string     `json:"password"`
	CreatedAt      time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
