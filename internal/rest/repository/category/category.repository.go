package repository

import "github.com/labovector/vecsys-api/entity"

type CategoryRepository interface {
	CreateCategory(category *entity.Category) (entity.Category, error)
	GetCategory(id string) (entity.Category, error)
	GetAllCategories(eventId ...string) ([]entity.Category, error)
	UpdateCategory(id string, category *entity.Category) (entity.Category, error)
	DeleteCategory(id string) error
}
