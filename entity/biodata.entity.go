package entity

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Biodata struct {
	Id            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ParticipantId string     `json:"participant_id"`
	Name          string     `json:"name,omitempty"`
	Gender        Gender     `json:"gender,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	Email         string     `json:"email,omitempty"`
	IdNumber      string     `json:"id_number,omitempty"`
	IdCardPicture string     `json:"id_card_picture,omitempty"`
	CreatedAt     time.Time  `gorm:"default:now();" json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
