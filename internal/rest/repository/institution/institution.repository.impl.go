package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type institutionRepositoryImpl struct {
	db *gorm.DB
}

// WithDB implements InstitutionRepository.
func (i *institutionRepositoryImpl) WithDB(db *gorm.DB) InstitutionRepository {
	return &institutionRepositoryImpl{
		db: db,
	}
}

// CreateInstitution implements InstitutionRepository.
func (i *institutionRepositoryImpl) CreateInstitution(institution *entity.Institution) (*entity.Institution, error) {
	db := i.db.Create(institution)
	return institution, db.Error
}

// DeleteInstitution implements InstitutionRepository.
func (i *institutionRepositoryImpl) DeleteInstitution(id string) error {
	err := i.db.Delete(&entity.Institution{}, "id = ?", id).Error
	return err
}

// GetAllInstitutions implements InstitutionRepository.
func (i *institutionRepositoryImpl) GetAllInstitutions(eventId ...string) ([]entity.Institution, error) {
	var institutions []entity.Institution
	if len(eventId) > 0 {
		err := i.db.Find(&institutions, "event_id = ?", eventId[0]).Error
		return institutions, err
	}
	err := i.db.Find(&institutions).Error
	return institutions, err
}

// GetInstitution implements InstitutionRepository.
func (i *institutionRepositoryImpl) GetInstitution(id string) (*entity.Institution, error) {
	var institution entity.Institution
	err := i.db.Preload("Participants").First(&institution, "id = ?", id).Error
	return &institution, err
}

// UpdateInstitution implements InstitutionRepository.
func (i *institutionRepositoryImpl) UpdateInstitution(id string, institution *entity.Institution) error {
	err := i.db.Model(&entity.Institution{}).Where("id = ?", id).Updates(institution).Error
	return err
}

func NewInstitutionRepositoryImpl(db *gorm.DB) InstitutionRepository {
	return &institutionRepositoryImpl{
		db: db,
	}
}
