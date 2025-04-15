package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type VoucherRepositoryImpl struct {
	db *gorm.DB
}

// CreateVoucher implements VoucherRepository.
func (v *VoucherRepositoryImpl) CreateVoucher(voucher *entity.Voucher) (entity.Voucher, error) {
	db := v.db.Create(voucher)
	return *voucher, db.Error
}

// DeleteVoucher implements VoucherRepository.
func (v *VoucherRepositoryImpl) DeleteVoucher(id string) error {
	err := v.db.Delete(&entity.Voucher{}, "id = ?", id).Error
	return err
}

// GetAllVoucher implements VoucherRepository.
func (v *VoucherRepositoryImpl) GetAllVoucher(eventId ...string) ([]entity.Voucher, error) {
	var vouchers []entity.Voucher

	if len(eventId) > 0 {
		err := v.db.Find(&vouchers, "event_id = ?", eventId[0]).Error
		return vouchers, err
	}
	err := v.db.Find(&vouchers).Error
	return vouchers, err
}

// GetVoucherByCode implements VoucherRepository.
func (v *VoucherRepositoryImpl) GetVoucherByCode(code string) (*entity.Voucher, error) {
	voucher := &entity.Voucher{}
	err := v.db.First(voucher, "voucher = ?", code).Error
	return voucher, err
}

// UpdateVoucher implements VoucherRepository.
func (v *VoucherRepositoryImpl) UpdateVoucher(id string, voucher *entity.Voucher) error {
	err := v.db.Model(&entity.Voucher{}).Where("id = ?", id).Updates(voucher).Error
	return err
}

func NewVoucherRepositoryImpl(db *gorm.DB) VoucherRepository {
	return &VoucherRepositoryImpl{
		db: db,
	}
}
