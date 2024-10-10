package repository

import "github.com/labovector/vecsys-api/entity"

type EventRepository interface {
	CreateEvent(event *entity.Event) (entity.Event, error)
	FindAllEvent() ([]entity.Event, error)
	FindAllActiveEvent() ([]entity.Event, error)
	FindEventById(id string) (*entity.Event, error)
	UpdateEvent(id string, event *entity.Event) error
	DeleteEvent(id string) error
}
