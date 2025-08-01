package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateParticipant(participant *entity.Participant) (entity.Participant, error)
	BulkAddParticipant(participants []entity.Participant) error
	FindAllParticipant(eventId ...string) ([]entity.Participant, error)
	FindParticipantById(id string, preload bool) (*entity.Participant, error)
	FindParticipantByEmail(email string) (*entity.Participant, error)
	UpdateParticipant(id string, participant *entity.Participant) error
	DeleteParticipant(id string) error

	// Additional Query
	FindAllPaidParticipant(eventId string) ([]entity.Participant, error)
	FindAllUnpaidParticipant(eventId string) ([]entity.Participant, error)

	// Biodata
	FindBiodataByParticipantId(participantId string) ([]entity.Biodata, error)
	FindBiodataById(id string) (*entity.Biodata, error)
	AddBiodata(participantId *string, biodata *entity.Biodata) (entity.Biodata, error)
	BulkAddBiodata(biodatas []entity.Biodata) error
	UpdateBiodata(id string, biodata *entity.Biodata) error
	BulkUpdateBiodata(participantId string, biodatas []entity.Biodata) error
	RemoveBiodata(id string) error

	// Custom tx wrapper returning custom paymentrepository
	WithDB(db *gorm.DB) UserRepository
}
