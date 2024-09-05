package repository

import (
	"github.com/labovector/vecsys-api/database"
	"github.com/labovector/vecsys-api/entity"
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
func (u *userRepositoryImpl) FindParticipantByEmail(email string) (*entity.Participant, error) {
	participant := &entity.Participant{}
	if err := database.DB.First(participant, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return participant, nil
}

// UpdateParticipant implements UserRepository.
func (u *userRepositoryImpl) UpdateParticipant(id string, participant *entity.Participant) error {
	db := database.DB.Model(&entity.Participant{}).Where("id = ?", id).Updates(participant)
	return db.Error
}

func NewUserRepositoryImpl() UserRepository {
	return &userRepositoryImpl{}
}
