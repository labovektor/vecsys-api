package repository

import "github.com/labovector/vecsys-api/entity"

type InstitutionRepository interface {
	CreateInstitution(institution *entity.Institution) (*entity.Institution, error)
	GetInstitution(id string) (*entity.Institution, error)
	GetAllInstitutions(eventId ...string) ([]entity.Institution, error)
	UpdateInstitution(id string, institution *entity.Institution) error
	DeleteInstitution(id string) error
}
