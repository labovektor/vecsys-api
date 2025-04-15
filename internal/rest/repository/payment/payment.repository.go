package repository

import "github.com/labovector/vecsys-api/entity"

type PaymentRepository interface {
	// Payment Options
	CreatePaymentOption(paymentOption *entity.PaymentOption) (entity.PaymentOption, error)
	UpdatePaymentOption(paymentOptionId string, paymentOption *entity.PaymentOption) (entity.PaymentOption, error)
	DeletePaymentOption(paymentOptionId string) (entity.PaymentOption, error)
	GetPaymentOptionById(id string) (entity.PaymentOption, error)
	GetPaymentOptions(eventId ...string) ([]entity.PaymentOption, error)

	// Payment
	CreatePayment(payment *entity.Payment) (entity.Payment, error)
	UpdatePayment(paymentId string, payment *entity.Payment) (entity.Payment, error)
	DeletePayment(paymentId string) (entity.Payment, error)
	GetPaymentById(id string) (entity.Payment, error)
	GetPaymentByParticipantId(participantId string) (entity.Payment, error)
	GetPayments(eventId ...string) ([]entity.Payment, error)
}
