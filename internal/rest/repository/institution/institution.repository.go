package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type InstitutionRepository interface {
	CreateInstitution(institution *entity.Institution) (*entity.Institution, error)
	GetInstitution(id string) (*entity.Institution, error)
	GetAllInstitutions(eventId ...string) ([]entity.Institution, error)
	UpdateInstitution(id string, institution *entity.Institution) error
	DeleteInstitution(id string) error

	// Custom tx wrapper returning custom paymentrepository
	WithDB(db *gorm.DB) InstitutionRepository
}
