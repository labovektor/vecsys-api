package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/institution"
	"github.com/labovector/vecsys-api/internal/util"
)

type InstitutionController struct {
	institutionRepo repository.InstitutionRepository
}

func NewInstitutionController(institutionRepo repository.InstitutionRepository) *InstitutionController {
	return &InstitutionController{
		institutionRepo: institutionRepo,
	}
}

func (cc *InstitutionController) GetAllInstitutions(c *fiber.Ctx) error {
	institutions, err := cc.institutionRepo.GetAllInstitutions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting institutions"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   institutions,
	})
}

func (cc *InstitutionController) GetInstitutions(c *fiber.Ctx) error {
	institutionId := c.Params("id")
	institution, err := cc.institutionRepo.GetInstitution(institutionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting institution"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   institution,
	})
}

func (cc *InstitutionController) AddInstitutions(c *fiber.Ctx) error {
	eventId := c.Params("id")
	req := new(dto.InstitutionAddReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	institution := entity.Institution{
		EventId:         &eventId,
		Name:            req.Name,
		Email:           req.Email,
		PendampingName:  req.PendampingName,
		PendampingPhone: req.PendampingPhone,
	}
	newInstitution, err := cc.institutionRepo.CreateInstitution(&institution)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding institution"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   newInstitution,
	})
}

func (cc *InstitutionController) UpdateInstitutions(c *fiber.Ctx) error {
	institutionId := c.Params("id")
	req := new(dto.InstitutionUpdateReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
		})
	}

	institution := entity.Institution{
		Name:            req.Name,
		Email:           req.Email,
		PendampingName:  req.PendampingName,
		PendampingPhone: req.PendampingPhone,
	}
	err := cc.institutionRepo.UpdateInstitution(institutionId, &institution)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating institution"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   "Institution Updated Successfully",
	})
}

func (cc *InstitutionController) DeleteInstitutions(c *fiber.Ctx) error {
	institutionId := c.Params("id")
	err := cc.institutionRepo.DeleteInstitution(institutionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when deleting institution"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   "Institution Deleted Successfully",
	})
}
