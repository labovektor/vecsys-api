package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentOption struct {
	Id        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	EventId   string     `json:"event_id"`
	Event     *Event     `json:"event" gorm:"foreignKey:EventId;references:Id"`
	Provider  string     `json:"provider"`
	Account   string     `json:"account"`
	Name      string     `json:"name"`
	AsQR      *bool      `gorm:"default:false" json:"as_qr"`
	CreatedAt time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
