package repository

import (
	"github.com/labovector/vecsys-api/database"
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type userRepositoryImpl struct{}

// CreateParticipant implements UserRepository.
func (u *userRepositoryImpl) CreateParticipant(participant *entity.Participant) (entity.Participant, error) {
	db := database.DB.Create(participant)
	return *participant, db.Error
}

// DeleteParticipant implements UserRepository.
func (u *userRepositoryImpl) DeleteParticipant(id string) error {
	db := database.DB.Where("id = ?", id).Delete(&entity.Participant{})
	return db.Error
}

// FindAllParticipant implements UserRepository.
func (u *userRepositoryImpl) FindAllParticipant() ([]entity.Participant, error) {
	var participants []entity.Participant
	if err := database.DB.Find(&participants).Error; err != nil {
		return nil, err
	}
	return participants, nil
}

// FindParticipantById implements UserRepository.
func (u *userRepositoryImpl) FindParticipantById(id string) (*entity.Participant, error) {
	participant := &entity.Participant{}
	if err := database.DB.First(participant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return participant, nil
}

// FindParticipantByUsername implements UserRepository.
func (u *userRepositoryImpl) FindParticipantByUsername(username string) (*entity.Participant, error) {
	participant := &entity.Participant{}
	if err := database.DB.First(participant, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return participant, nil
}

// IsParticipantExist implements UserRepository.
func (u *userRepositoryImpl) IsParticipantExist(username string) (bool, error) {
	participant := &entity.Participant{}
	if err := database.DB.First(participant, "username = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User does not exist
			return false, nil
		}
		// Other errors (e.g., DB connection issues)
		return false, err
	}

	// User exists
	return true, nil
}

// UpdateParticipant implements UserRepository.
func (u *userRepositoryImpl) UpdateParticipant(id string, participant *entity.Participant) error {
	db := database.DB.Model(&entity.Participant{}).Where("id = ?", id).Updates(participant)
	return db.Error
}

func NewUserRepositoryImpl() UserRepository {
	return &userRepositoryImpl{}
}
