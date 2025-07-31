package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ParticipantId   *string        `json:"participant_id"`
	BankName        string         `json:"bank_name"`
	BankAccount     string         `json:"bank_account"`
	ReferalId       *string        `json:"referal_id"`
	Referal         *Referal       `json:"referal" gorm:"foreignKey:ReferalId;references:Id"`
	Date            *time.Time     `json:"date"`
	Invoice         string         `json:"invoice"`
	PaymentOptionId *string        `json:"payment_option_id"`
	PaymentOption   *PaymentOption `json:"payment_option" gorm:"foreignKey:PaymentOptionId;references:Id"`
	CreatedAt       time.Time      `gorm:"default:now();" json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
}
