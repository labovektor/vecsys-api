package repository

import (
	"github.com/labovector/vecsys-api/database"
	"github.com/labovector/vecsys-api/entity"
)

type userRepositoryImpl struct{}

// FindBiodataById implements UserRepository.
func (u *userRepositoryImpl) FindBiodataById(id string) (*entity.Biodata, error) {
	biodata := &entity.Biodata{}
	if err := database.DB.First(biodata, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return biodata, nil
}

// FindBiodataByParticipantId implements UserRepository.
func (u *userRepositoryImpl) FindBiodataByParticipantId(participantId string) ([]entity.Biodata, error) {
	var biodatas []entity.Biodata
	if err := database.DB.Find(&biodatas, "participant_id = ?", participantId).Error; err != nil {
		return nil, err
	}
	return biodatas, nil
}

// AddBiodata implements UserRepository.
func (u *userRepositoryImpl) AddBiodata(participantId string, biodata *entity.Biodata) (entity.Biodata, error) {
	biodata.ParticipantId = participantId

	db := database.DB.Create(biodata)
	return *biodata, db.Error
}

// RemoveBiodata implements UserRepository.
func (u *userRepositoryImpl) RemoveBiodata(id string) error {
	db := database.DB.Where("id = ?", id).Delete(&entity.Biodata{})
	return db.Error
}

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
