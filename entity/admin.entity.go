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

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpReq struct {
	Username       string `json:"username" form:"username"`
	DisplayName    string `json:"display_name" form:"display_name"`
	ProfilePicture string `json:"profile_picture"`
	Email          string `json:"email" form:"email"`
	Password       string `json:"password" form:"password"`
}
