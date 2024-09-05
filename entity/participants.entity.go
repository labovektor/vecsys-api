package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	Biodata       BiodataSlice `json:"biodata"`
	PaymentDataId string       `gorm:"unique" json:"payment_data_id"`
	Payment       *Payment     `json:"payment_data" gorm:"foreignKey:PaymentDataId;references:Id"`
	VerifiedAt    *time.Time   `json:"verified_at,omitempty"`
	CreatedAt     time.Time    `gorm:"default:now();" json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
}

type Biodata struct {
	Name          string `json:"name,omitempty"`
	Gender        string `json:"gender,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	IdNumber      string `json:"id_number,omitempty"`
	IdCardPicture string `json:"id_card_picture,omitempty"`
}

// Define a custom type for the slice of Biodata
type BiodataSlice []Biodata

// Implement the Valuer interface for BiodataSlice
func (b *BiodataSlice) Value() (driver.Value, error) {
	// Convert the slice of Biodata to JSON
	if b == nil {
		return nil, nil
	}
	return json.Marshal(b)
}

// Implement the Scanner interface for BiodataSlice
func (b *BiodataSlice) Scan(src interface{}) error {
	if src == nil {
		*b = nil
		return nil
	}
	switch data := src.(type) {
	case []byte:
		// Convert the JSON data to a slice of Biodata
		return json.Unmarshal(data, b)
	case string:
		// Convert the JSON string data to a slice of Biodata
		return json.Unmarshal([]byte(data), b)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
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
