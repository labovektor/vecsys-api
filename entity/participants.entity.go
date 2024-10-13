package entity

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	Id            uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	EventId       string       `json:"event_id"`
	Event         *Event       `json:"event" gorm:"foreignKey:EventId;references:Id"`
	RegionId      string       `json:"region_id"`
	Region        *Region      `json:"region" gorm:"foreignKey:RegionId;references:Id"`
	CategoryId    string       `json:"category_id"`
	Category      *Category    `json:"category" gorm:"foreignKey:CategoryId;references:Id"`
	Name          string       `json:"name"`
	InstitutionId string       `json:"institution_id"`
	Institution   *Institution `json:"institution" gorm:"foreignKey:InstitutionId;references:Id"`
	Email         string       `gorm:"unique" json:"email"`
	Password      string       `json:"password"`
	PaymentDataId string       `gorm:"unique" json:"payment_data_id"`
	Payment       *Payment     `json:"payment_data" gorm:"foreignKey:PaymentDataId;references:Id"`
	VerifiedAt    *time.Time   `json:"verified_at,omitempty"`
	LockedAt      *time.Time   `json:"locked_at,omitempty"`
	CreatedAt     time.Time    `gorm:"default:now();" json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
}

type ParticipantLoginReq struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type ParticipantSignUpReq struct {
	EventId  string `json:"event_id" form:"event_id"`
	Name     string `json:"name" form:"name"`
	Email    string ` json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
