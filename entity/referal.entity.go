package entity

import (
	"time"

	"github.com/google/uuid"
)

type Referal struct {
	Id            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventId       string     `json:"event_id"`
	Event         *Event     `json:"event" gorm:"foreignKey:EventId;references:Id"`
	Code          string     `json:"code"`
	Desc          string     `json:"desc"`
	IsDiscount    *bool      `gorm:"default:false" json:"is_voucher"`
	Discount      int        `gorm:"default:0" json:"discount"`
	SeatAvailable int        `gorm:"default:1" json:"seat_available"`
	CreatedAt     time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
