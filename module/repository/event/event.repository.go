package repository

import "github.com/labovector/vecsys-api/entity"

type EventRepository interface {
	CreateEvent(event *entity.Event) (entity.Event, error)
	FindAllEvent(adminId ...string) ([]entity.Event, error)
	FindAllActiveEvent(adminId ...string) ([]entity.Event, error)
	FindEventById(id string, adminId ...string) (*entity.Event, error)
	UpdateEvent(id string, event *entity.Event) error
	DeleteEvent(id string) error
}
