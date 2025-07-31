package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type ReferalRepository interface {
	GetAllVoucher(eventId ...string) ([]entity.Referal, error)
	GetAllNonDiscountVoucher(eventId ...string) ([]entity.Referal, error)
	GetAllDiscountVoucher(eventId ...string) ([]entity.Referal, error)
	GetVoucherByCode(code string, eventId ...string) (*entity.Referal, error)
	CreateVoucher(voucher *entity.Referal) (entity.Referal, error)
	UpdateVoucher(id string, voucher *entity.Referal) error
	DeleteVoucher(id string) error

	// Custom tx wrapper returning custom paymentrepository
	WithDB(db *gorm.DB) ReferalRepository
}
