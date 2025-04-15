package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type paymentRepositoryImpl struct {
	db *gorm.DB
}

// Payment Options

// CreatePaymentOption implements PaymentRepository.
func (p *paymentRepositoryImpl) CreatePaymentOption(paymentOption *entity.PaymentOption) (entity.PaymentOption, error) {
	db := p.db.Create(paymentOption)
	return *paymentOption, db.Error
}

// DeletePaymentOption implements PaymentRepository.
func (p *paymentRepositoryImpl) DeletePaymentOption(paymentOptionId string) error {
	err := p.db.Delete(&entity.PaymentOption{}, "id = ?", paymentOptionId).Error
	return err
}

// GetPaymentOptionById implements PaymentRepository.
func (p *paymentRepositoryImpl) GetPaymentOptionById(id string) (entity.PaymentOption, error) {
	var paymentOption entity.PaymentOption
	err := p.db.First(&paymentOption, "id = ?", id).Error
	return paymentOption, err
}

// GetPaymentOptions implements PaymentRepository.
func (p *paymentRepositoryImpl) GetPaymentOptions(eventId ...string) ([]entity.PaymentOption, error) {
	var paymentOptions []entity.PaymentOption
	if len(eventId) > 0 {
		err := p.db.Find(&paymentOptions, "event_id = ?", eventId[0]).Error
		return paymentOptions, err
	}
	err := p.db.Find(&paymentOptions).Error
	return paymentOptions, err
}

func (p *paymentRepositoryImpl) UpdatePaymentOption(paymentOptionId string, paymentOption *entity.PaymentOption) error {
	err := p.db.Model(&entity.PaymentOption{}).Where("id = ?", paymentOptionId).Updates(paymentOption).Error
	return err
}

// =========================================================================

// Payment

// CreatePayment implements PaymentRepository.
func (p *paymentRepositoryImpl) CreatePayment(payment *entity.Payment) (entity.Payment, error) {
	db := p.db.Create(payment)
	return *payment, db.Error
}

// DeletePayment implements PaymentRepository.
func (p *paymentRepositoryImpl) DeletePayment(paymentId string) error {
	err := p.db.Delete(&entity.Payment{}, "id = ?", paymentId).Error
	return err
}

// GetPaymentById implements PaymentRepository.
func (p *paymentRepositoryImpl) GetPaymentById(id string) (entity.Payment, error) {
	var payment entity.Payment
	err := p.db.First(&payment, "id = ?", id).Error
	return payment, err
}

// GetPaymentByParticipantId implements PaymentRepository.
func (p *paymentRepositoryImpl) GetPaymentByParticipantId(participantId string) (entity.Payment, error) {
	var payment entity.Payment
	err := p.db.First(&payment, "participant_id = ?", participantId).Error
	return payment, err
}

// GetPayments implements PaymentRepository.
func (p *paymentRepositoryImpl) GetPayments(eventId ...string) ([]entity.Payment, error) {
	var payments []entity.Payment
	if len(eventId) > 0 {
		err := p.db.Find(&payments, "event_id = ?", eventId[0]).Error
		return payments, err
	}
	err := p.db.Find(&payments).Error
	return payments, err
}

// UpdatePayment implements PaymentRepository.
func (p *paymentRepositoryImpl) UpdatePayment(paymentId string, payment *entity.Payment) error {
	err := p.db.Model(&entity.Payment{}).Where("id = ?", paymentId).Updates(payment).Error
	return err
}

// UpdatePaymentOption implements PaymentRepository.

func NewPaymentRepositoryImpl(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{
		db: db,
	}
}
