package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type ReferalRepositoryImpl struct {
	db *gorm.DB
}

// CreateVoucher implements VoucherRepository.
func (v *ReferalRepositoryImpl) CreateVoucher(voucher *entity.Referal) (entity.Referal, error) {
	db := v.db.Create(voucher)
	return *voucher, db.Error
}

// DeleteVoucher implements VoucherRepository.
func (v *ReferalRepositoryImpl) DeleteVoucher(id string) error {
	err := v.db.Delete(&entity.Referal{}, "id = ?", id).Error
	return err
}

// GetAllVoucher implements VoucherRepository.
func (v *ReferalRepositoryImpl) GetAllVoucher(eventId ...string) ([]entity.Referal, error) {
	var vouchers []entity.Referal

	if len(eventId) > 0 {
		err := v.db.Find(&vouchers, "event_id = ?", eventId[0]).Error
		return vouchers, err
	}
	err := v.db.Find(&vouchers).Error
	return vouchers, err
}

// GetVoucherByCode implements VoucherRepository.
func (v *ReferalRepositoryImpl) GetVoucherByCode(code string) (*entity.Referal, error) {
	voucher := &entity.Referal{}
	err := v.db.First(voucher, "voucher = ?", code).Error
	return voucher, err
}

// UpdateVoucher implements VoucherRepository.
func (v *ReferalRepositoryImpl) UpdateVoucher(id string, voucher *entity.Referal) error {
	err := v.db.Model(&entity.Referal{}).Where("id = ?", id).Updates(voucher).Error
	return err
}

func NewReferalRepositoryImpl(db *gorm.DB) ReferalRepository {
	return &ReferalRepositoryImpl{
		db: db,
	}
}
