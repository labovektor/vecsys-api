package repository

import (
	"github.com/labovector/vecsys-api/database"
	"github.com/labovector/vecsys-api/entity"
)

type adminRepositoryImpl struct{}

// CreateAdmin implements AdminRepository.
func (a adminRepositoryImpl) CreateAdmin(admin *entity.Admin) (entity.Admin, error) {
	db := database.DB.Create(admin)
	return *admin, db.Error
}

// DeleteAdmin implements AdminRepository.
func (a adminRepositoryImpl) DeleteAdmin(id string) error {
	db := database.DB.Where("id = ?", id).Delete(&entity.Admin{})
	return db.Error
}

// FindAdminById implements AdminRepository.
func (a adminRepositoryImpl) FindAdminById(id string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	if err := database.DB.First(admin, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// FindAdminByUsername implements AdminRepository.
func (a adminRepositoryImpl) FindAdminByUsername(username string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	if err := database.DB.First(admin, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// FindAllAdmin implements AdminRepository.
func (a adminRepositoryImpl) FindAllAdmin() ([]entity.Admin, error) {
	var admins []entity.Admin
	if err := database.DB.Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// UpdateAdmin implements AdminRepository.
func (a adminRepositoryImpl) UpdateAdmin(id string, admin *entity.Admin) error {
	db := database.DB.Model(&entity.Admin{}).Where("id = ?", id).Updates(admin)
	return db.Error
}

func NewAdminRepositoryImpl() AdminRepository {
	return adminRepositoryImpl{}
}
