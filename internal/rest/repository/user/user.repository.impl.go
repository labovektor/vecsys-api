package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// WithDB implements UserRepository.
func (u *userRepositoryImpl) WithDB(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// BulkAddParticipant implements UserRepository.
func (u *userRepositoryImpl) BulkAddParticipant(participants []entity.Participant) error {
	tx := u.db.Begin()
	for _, participant := range participants {
		db := tx.Create(&participant)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}
	}
	return tx.Commit().Error
}

// BulkUpdateBiodata implements UserRepository.
func (u *userRepositoryImpl) BulkUpdateBiodata(participantId string, biodatas []entity.Biodata) error {
	tx := u.db.Begin()
	for _, biodata := range biodatas {
		db := tx.Model(&entity.Biodata{}).Where("id = ?", biodata.Id.String()).Where("participant_id = ?", participantId).Updates(biodata)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}
	}
	return tx.Commit().Error
}

// BulkAddBiodata implements UserRepository.
func (u *userRepositoryImpl) BulkAddBiodata(biodatas []entity.Biodata) error {
	tx := u.db.Begin()
	for _, biodata := range biodatas {
		db := tx.Create(&biodata)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}
	}
	return tx.Commit().Error
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
func (u *userRepositoryImpl) AddBiodata(participantId *string, biodata *entity.Biodata) (entity.Biodata, error) {
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
func (u *userRepositoryImpl) FindAllParticipant(eventId ...string) ([]entity.Participant, error) {
	var participants []entity.Participant
	if len(eventId) > 0 {
		err := u.db.Find(&participants, "event_id = ?", eventId[0]).Error
		return participants, err
	}
	err := u.db.Find(&participants).Error
	return participants, err
}

// FindAllPaidParticipant implements UserRepository.
func (u *userRepositoryImpl) FindAllPaidParticipant(eventId string) ([]entity.Participant, error) {
	var participants []entity.Participant

	err := u.db.Preload("Payment").Preload("Region").Preload("Category").
		Where("event_id = ?", eventId).
		Where("progress_step IN ?", []entity.ParticipantProgress{
			entity.StepPaidParticipant,
			entity.StepVerifiedParticipant,
			entity.StepSelectInstitutionParticipant,
			entity.StepFillBiodatasParticipant,
			entity.StepLockedParticipant,
		}).Find(&participants).Error
	return participants, err
}

// FindAllUnpaidParticipant implements UserRepository.
func (u *userRepositoryImpl) FindAllUnpaidParticipant(eventId string) ([]entity.Participant, error) {
	var participants []entity.Participant

	err := u.db.Preload("Payment").Preload("Region").Preload("Category").
		Where("event_id = ?", eventId).
		Where("progress_step IN ?", []entity.ParticipantProgress{
			entity.StepRegisteredParticipant,
			entity.StepCategorizedParticipant,
		}).Find(&participants).Error
	return participants, err
}

// FindParticipantById implements UserRepository.
func (u *userRepositoryImpl) FindParticipantById(id string, preload bool) (*entity.Participant, error) {
	participant := &entity.Participant{}
	if preload {
		if err := u.db.Preload("Event").Preload("Biodata").Preload("Institution").Preload("Region").Preload("Category").Preload("Payment").First(participant, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return participant, nil
	}
	if err := u.db.Preload("Event").First(participant, "id = ?", id).Error; err != nil {
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
