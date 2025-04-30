package entity

import (
	"time"

	"github.com/google/uuid"
)

type Region struct {
	Id            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventId       *string    `json:"event_id"`
	Event         *Event     `json:"event" gorm:"foreignKey:EventId;references:Id"`
	Name          string     `json:"name"`
	Visible       *bool      `gorm:"default:true" json:"visible"`
	ContactNumber string     `json:"contact_number"`
	ContactName   string     `json:"contact_name"`
	CreatedAt     time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
