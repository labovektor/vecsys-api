package repository

import "github.com/labovector/vecsys-api/entity"

type UserRepository interface {
	CreateParticipant(participant *entity.Participant) (entity.Participant, error)
	FindAllParticipant() ([]entity.Participant, error)
	FindParticipantById(id string) (*entity.Participant, error)
	FindParticipantByUsername(username string) (*entity.Participant, error)
	IsParticipantExist(username string) (bool, error)
	UpdateParticipant(id string, participant *entity.Participant) error
	DeleteParticipant(id string) error
}
