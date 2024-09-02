package entity

import (
	"time"

	"github.com/google/uuid"
)

type Institution struct {
	Id              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	EventId         string     `json:"event_id"`
	Event           *Event     `json:"event" gorm:"foreignKey:EventId;references:Id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	PendampingName  string     `json:"pendamping_name"`
	PendampingPhone string     `json:"pendamping_phone"`
	CreatedAt       time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
