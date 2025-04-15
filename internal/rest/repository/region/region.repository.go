package repository

import "github.com/labovector/vecsys-api/entity"

type RegionRepository interface {
	GetRegion(id string) (entity.Region, error)
	GetAllRegion(eventId ...string) ([]entity.Region, error)
	CreateRegion(region entity.Region) (entity.Region, error)
	UpdateRegion(id string, region *entity.Region) error
	DeleteRegion(id string) error
}
