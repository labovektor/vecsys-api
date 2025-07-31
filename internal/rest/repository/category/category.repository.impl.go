package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepositoryImpl) CreateCategory(category *entity.Category) (entity.Category, error) {
	db := c.db.Create(category)
	return *category, db.Error
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepositoryImpl) DeleteCategory(id string) error {
	err := c.db.Delete(&entity.Category{}, "id = ?", id).Error
	return err
}

// GetAllCategories implements CategoryRepository.
func (c *categoryRepositoryImpl) GetAllCategories(eventId ...string) ([]entity.Category, error) {
	var categories []entity.Category
	if len(eventId) > 0 {
		err := c.db.Find(&categories, "event_id = ?", eventId[0]).Error
		return categories, err
	}
	err := c.db.Find(&categories).Error
	return categories, err
}

// GetAllActiveCategories implements CategoryRepository.
func (c *categoryRepositoryImpl) GetAllActiveCategories(eventId ...string) ([]entity.Category, error) {
	var categories []entity.Category
	if len(eventId) > 0 {
		err := c.db.Find(&categories, "event_id = ? AND visible = ?", eventId[0], true).Error
		return categories, err
	}
	err := c.db.Find(&categories, "visible = ?", true).Error
	return categories, err
}

// GetCategory implements CategoryRepository.
func (c *categoryRepositoryImpl) GetCategory(id string) (entity.Category, error) {
	var category entity.Category
	err := c.db.First(&category, "id = ?", id).Error
	return category, err
}

// UpdateCategory implements CategoryRepository.
func (c *categoryRepositoryImpl) UpdateCategory(id string, category *entity.Category) (entity.Category, error) {
	err := c.db.Model(&entity.Category{}).Where("id = ?", id).Updates(category).Error
	return *category, err
}

func NewCategoryRepositoryImpl(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}
