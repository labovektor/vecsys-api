package repository

import "github.com/labovector/vecsys-api/entity"

type AdminRepository interface {
	CreateAdmin(admin *entity.Admin) (entity.Admin, error)
	FindAllAdmin() ([]entity.Admin, error)
	FindAdminById(id string) (*entity.Admin, error)
	FindAdminByUsername(username string) (*entity.Admin, error)
	UpdateAdmin(id string, admin *entity.Admin) error
	DeleteAdmin(id string) error
}
