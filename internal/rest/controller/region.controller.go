package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/region"
	"github.com/labovector/vecsys-api/internal/util"
)

type RegionController struct {
	regionRepo repository.RegionRepository
}

func NewRegionController(regionRepo repository.RegionRepository) RegionController {
	return RegionController{regionRepo: regionRepo}
}

func (rc *RegionController) GetAllRegionsByEventId(c *fiber.Ctx) error {
	eventId := c.Params("id")
	regions, err := rc.regionRepo.GetAllRegion(eventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting regions"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   regions,
	})
}

func (rc *RegionController) GetRegionById(c *fiber.Ctx) error {
	regionId := c.Params("id")
	region, err := rc.regionRepo.GetRegion(regionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting region"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   region,
	})
}

func (rc *RegionController) AddRegionToEvent(c *fiber.Ctx) error {
	eventId := c.Params("id")
	req := new(dto.RegionCreateReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	region := entity.Region{
		EventId:       &eventId,
		Name:          req.Name,
		ContactNumber: req.ContactNumber,
		ContactName:   req.ContactName,
		Visible:       req.Visible,
	}

	region, err := rc.regionRepo.CreateRegion(region)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding region"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   region,
	})
}

func (rc *RegionController) UpdateRegion(c *fiber.Ctx) error {
	regionId := c.Params("id")
	req := new(dto.RegionEditReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	region := entity.Region{
		Name:          req.Name,
		ContactNumber: req.ContactNumber,
		ContactName:   req.ContactName,
		Visible:       req.Visible,
	}

	err := rc.regionRepo.UpdateRegion(regionId, &region)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating region"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   "Region Updated Successfully",
	})
}

func (rc *RegionController) DeleteRegion(c *fiber.Ctx) error {
	regionId := c.Params("id")
	err := rc.regionRepo.DeleteRegion(regionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when deleting region"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   "Region Deleted Successfully",
	})
}
