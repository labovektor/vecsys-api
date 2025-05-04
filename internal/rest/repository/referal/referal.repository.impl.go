package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type ReferalRepositoryImpl struct {
	db *gorm.DB
}

// GetAllDiscountVoucher implements ReferalRepository.
func (v *ReferalRepositoryImpl) GetAllDiscountVoucher(eventId ...string) ([]entity.Referal, error) {
	var vouchers []entity.Referal

	if len(eventId) > 0 {
		err := v.db.Find(&vouchers, "event_id = ? AND is_discount = ?", eventId[0], true).Error
		return vouchers, err
	}
	err := v.db.Find(&vouchers, "is_discount = ?", true).Error
	return vouchers, err
}

// GetAllNonDiscountVoucher implements ReferalRepository.
func (v *ReferalRepositoryImpl) GetAllNonDiscountVoucher(eventId ...string) ([]entity.Referal, error) {
	var vouchers []entity.Referal

	if len(eventId) > 0 {
		err := v.db.Find(&vouchers, "event_id = ? AND is_discount = ?", eventId[0], false).Error
		return vouchers, err
	}
	err := v.db.Find(&vouchers, "is_discount = ?", false).Error
	return vouchers, err
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
func (v *ReferalRepositoryImpl) GetVoucherByCode(code string, eventId ...string) (*entity.Referal, error) {
	voucher := &entity.Referal{}
	if len(eventId) > 0 {
		err := v.db.First(voucher, "code = ? AND event_id = ?", code, eventId[0]).Error
		return voucher, err
	}
	err := v.db.First(voucher, "code = ?", code).Error
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
