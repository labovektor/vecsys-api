package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Id        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	EventId   string     `json:"event_id"`
	Event     *Event     `json:"event" gorm:"foreignKey:EventId;references:Id"`
	Name      string     `json:"name"`
	IsGroup   *bool      `gorm:"default:false" json:"is_group"`
	Visible   *bool      `gorm:"default:true" json:"visible"`
	CreatedAt time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
