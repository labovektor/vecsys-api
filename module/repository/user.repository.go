package repository

import "github.com/labovector/vecsys-api/entity"

type UserRepository interface {
	CreateParticipant(participant *entity.Participant) (entity.Participant, error)
	FindAllParticipant() ([]entity.Participant, error)
	FindParticipantById(id string) (*entity.Participant, error)
	FindParticipantByEmail(email string) (*entity.Participant, error)
	UpdateParticipant(id string, participant *entity.Participant) error
	DeleteParticipant(id string) error

	// Biodata
	FindBiodataByParticipantId(participantId string) ([]entity.Biodata, error)
	FindBiodataById(id string) (*entity.Biodata, error)
	AddBiodata(participantId string, biodata *entity.Biodata) (entity.Biodata, error)
	RemoveBiodata(id string) error
}
