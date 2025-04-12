package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type adminRepositoryImpl struct {
	db *gorm.DB
}

// CreateAdmin implements AdminRepository.
func (a adminRepositoryImpl) CreateAdmin(admin *entity.Admin) (entity.Admin, error) {
	return *admin, a.db.Create(admin).Error
}

// DeleteAdmin implements AdminRepository.
func (a adminRepositoryImpl) DeleteAdmin(id string) error {
	db := a.db.Where("id = ?", id).Delete(&entity.Admin{})
	return db.Error
}

// FindAdminById implements AdminRepository.
func (a adminRepositoryImpl) FindAdminById(id string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	if err := a.db.First(admin, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// FindAdminByUsername implements AdminRepository.
func (a adminRepositoryImpl) FindAdminByUsername(username string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	if err := a.db.First(admin, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// FindAllAdmin implements AdminRepository.
func (a adminRepositoryImpl) FindAllAdmin() ([]entity.Admin, error) {
	var admins []entity.Admin
	if err := a.db.Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// UpdateAdmin implements AdminRepository.
func (a adminRepositoryImpl) UpdateAdmin(id string, admin *entity.Admin) error {
	db := a.db.Model(&entity.Admin{}).Where("id = ?", id).Updates(admin)
	return db.Error
}

func NewAdminRepositoryImpl(db *gorm.DB) AdminRepository {
	return adminRepositoryImpl{
		db: db,
	}
}
