package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type eventRepositoryImpl struct {
	db *gorm.DB
}

// CreateEvent implements EventRepository.
func (e *eventRepositoryImpl) CreateEvent(event *entity.Event) (entity.Event, error) {
	db := e.db.Create(event)
	return *event, db.Error
}

// DeleteEvent implements EventRepository.
func (e *eventRepositoryImpl) DeleteEvent(id string) error {
	db := e.db.Where("id = ?", id).Delete(&entity.Event{})
	return db.Error
}

// FindAllActiveEvent implements EventRepository.
func (e *eventRepositoryImpl) FindAllActiveEvent(adminId ...string) ([]entity.Event, error) {
	var events []entity.Event
	if len(adminId) > 0 {
		if err := e.db.Where("active = ? AND admin_id = ?", true, adminId[0]).Find(&events).Error; err != nil {
			return nil, err
		}
	} else {
		if err := e.db.Where("active = ?", true).Find(&events).Error; err != nil {
			return nil, err
		}
	}
	return events, nil
}

// FindAllEvent implements EventRepository.
func (e *eventRepositoryImpl) FindAllEvent(adminId ...string) ([]entity.Event, error) {
	var events []entity.Event
	if len(adminId) > 0 {
		if err := e.db.Find(&events, "admin_id = ?", adminId[0]).Error; err != nil {
			return nil, err
		}
	} else {
		if err := e.db.Find(&events).Error; err != nil {
			return nil, err
		}
	}
	return events, nil
}

// FindEventById implements EventRepository.
func (e *eventRepositoryImpl) FindEventById(id string, adminId ...string) (*entity.Event, error) {
	event := &entity.Event{}
	if len(adminId) > 0 {
		if err := e.db.First(event, "id = ? AND admin_id = ?", id, adminId[0]).Error; err != nil {
			return nil, err
		}
	} else {
		if err := e.db.First(event, "id = ?", id).Error; err != nil {
			return nil, err
		}
	}
	return event, nil
}

// UpdateEvent implements EventRepository.
func (e *eventRepositoryImpl) UpdateEvent(id string, event *entity.Event) error {
	db := e.db.Model(&entity.Event{}).Where("id = ?", id).Updates(event)
	return db.Error
}

func NewEventRepositositoryImpl(db *gorm.DB) EventRepository {
	return &eventRepositoryImpl{
		db: db,
	}
}
