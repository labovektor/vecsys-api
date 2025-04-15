package repository

import (
	"github.com/labovector/vecsys-api/entity"
	"gorm.io/gorm"
)

type RegionRepositoryImpl struct {
	db *gorm.DB
}

// CreateRegion implements RegionRepository.
func (r *RegionRepositoryImpl) CreateRegion(region entity.Region) (entity.Region, error) {
	db := r.db.Create(&region)
	return region, db.Error
}

// DeleteRegion implements RegionRepository.
func (r *RegionRepositoryImpl) DeleteRegion(id string) error {
	err := r.db.Delete(&entity.Region{}, "id = ?", id).Error
	return err
}

// GetAllRegion implements RegionRepository.
func (r *RegionRepositoryImpl) GetAllRegion(eventId ...string) ([]entity.Region, error) {
	var regions []entity.Region
	if len(eventId) > 0 {
		err := r.db.Find(&regions, "event_id = ?", eventId[0]).Error
		return regions, err
	}
	err := r.db.Find(&regions).Error
	return regions, err
}

// GetRegion implements RegionRepository.
func (r *RegionRepositoryImpl) GetRegion(id string) (entity.Region, error) {
	var region entity.Region
	err := r.db.First(&region, "id = ?", id).Error
	return region, err
}

// UpdateRegion implements RegionRepository.
func (r *RegionRepositoryImpl) UpdateRegion(id string, region *entity.Region) error {
	err := r.db.Model(&entity.Region{}).Where("id = ?", id).Updates(region).Error
	return err
}

func NewRegionRepositoryImpl(db *gorm.DB) RegionRepository {
	return &RegionRepositoryImpl{
		db: db,
	}
}
