package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// UpdateBiodata implements UserRepository.
func (u *userRepositoryImpl) UpdateBiodata(id string, biodata *entity.Biodata) error {
	err := u.db.Model(&entity.Biodata{}).Where("id = ?", id).Updates(biodata).Error
	return err
}

// FindBiodataById implements UserRepository.
func (u *userRepositoryImpl) FindBiodataById(id string) (*entity.Biodata, error) {
	biodata := &entity.Biodata{}
	if err := u.db.First(biodata, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return biodata, nil
}

// FindBiodataByParticipantId implements UserRepository.
func (u *userRepositoryImpl) FindBiodataByParticipantId(participantId string) ([]entity.Biodata, error) {
	var biodatas []entity.Biodata
	if err := u.db.Find(&biodatas, "participant_id = ?", participantId).Error; err != nil {
		return nil, err
	}
	return biodatas, nil
}

// AddBiodata implements UserRepository.
func (u *userRepositoryImpl) AddBiodata(participantId string, biodata *entity.Biodata) (entity.Biodata, error) {
	biodata.ParticipantId = participantId

	db := u.db.Create(biodata)
	return *biodata, db.Error
}

// RemoveBiodata implements UserRepository.
func (u *userRepositoryImpl) RemoveBiodata(id string) error {
	db := u.db.Where("id = ?", id).Delete(&entity.Biodata{})
	return db.Error
}

// CreateParticipant implements UserRepository.
func (u *userRepositoryImpl) CreateParticipant(participant *entity.Participant) (entity.Participant, error) {
	db := u.db.Create(participant)
	return *participant, db.Error
}

// DeleteParticipant implements UserRepository.
func (u *userRepositoryImpl) DeleteParticipant(id string) error {
	db := u.db.Where("id = ?", id).Delete(&entity.Participant{})
	return db.Error
}

// FindAllParticipant implements UserRepository.
func (u *userRepositoryImpl) FindAllParticipant() ([]entity.Participant, error) {
	var participants []entity.Participant
	if err := u.db.Find(&participants).Error; err != nil {
		return nil, err
	}
	return participants, nil
}

// FindParticipantById implements UserRepository.
func (u *userRepositoryImpl) FindParticipantById(id string) (*entity.Participant, error) {
	participant := &entity.Participant{}
	if err := u.db.First(participant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return participant, nil
}

// FindParticipantByUsername implements UserRepository.
func (u *userRepositoryImpl) FindParticipantByEmail(email string) (*entity.Participant, error) {
	participant := &entity.Participant{}
	if err := u.db.First(participant, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return participant, nil
}

// UpdateParticipant implements UserRepository.
func (u *userRepositoryImpl) UpdateParticipant(id string, participant *entity.Participant) error {
	db := u.db.Model(&entity.Participant{}).Where("id = ?", id).Updates(participant)
	return db.Error
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}
