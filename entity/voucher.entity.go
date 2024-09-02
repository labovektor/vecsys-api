package entity

import (
	"time"

	"github.com/google/uuid"
)

type Voucher struct {
	Id            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	EventId       string     `json:"event_id"`
	Event         *Event     `json:"event" gorm:"foreignKey:EventId;references:Id"`
	Voucher       string     `json:"voucher"`
	Desc          string     `json:"desc"`
	SeatAvailable int        `gorm:"default:1" json:"seat_available"`
	CreatedAt     time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
