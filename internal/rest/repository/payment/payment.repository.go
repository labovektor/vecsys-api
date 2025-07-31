package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	// Payment Options
	CreatePaymentOption(paymentOption *entity.PaymentOption) (entity.PaymentOption, error)
	UpdatePaymentOption(paymentOptionId string, paymentOption *entity.PaymentOption) error
	DeletePaymentOption(paymentOptionId string) error
	GetPaymentOptionById(id string) (entity.PaymentOption, error)
	GetPaymentOptions(eventId ...string) ([]entity.PaymentOption, error)

	// Payment
	CreatePayment(payment *entity.Payment) (entity.Payment, error)
	UpdatePayment(paymentId string, payment *entity.Payment) error
	DeletePayment(paymentId string) error
	GetPaymentById(id string) (entity.Payment, error)
	GetPaymentByParticipantId(participantId string) (entity.Payment, error)
	GetPayments(eventId ...string) ([]entity.Payment, error)

	// Custom tx wrapper returning custom paymentrepository
	WithDB(db *gorm.DB) PaymentRepository
}
