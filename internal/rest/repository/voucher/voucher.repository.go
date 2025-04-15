package repository

import "github.com/labovector/vecsys-api/entity"

type VoucherRepository interface {
	GetAllVoucher(eventId ...string) ([]entity.Voucher, error)
	GetVoucherByCode(code string) (*entity.Voucher, error)
	CreateVoucher(voucher *entity.Voucher) (entity.Voucher, error)
	UpdateVoucher(id string, voucher *entity.Voucher) error
	DeleteVoucher(id string) error
}
