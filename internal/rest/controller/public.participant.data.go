package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	ir "github.com/labovector/vecsys-api/internal/rest/repository/institution"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
)

type ParticipantDataController struct {
	ParticipantRepository ur.UserRepository
	InstitutionRepository ir.InstitutionRepository
}

func NewParticipantDataController(
	participantRepository ur.UserRepository,
	institutionRepository ir.InstitutionRepository,
) *ParticipantDataController {
	return &ParticipantDataController{
		ParticipantRepository: participantRepository,
		InstitutionRepository: institutionRepository,
	}
}

func (p *ParticipantDataController) GetAllInstitution(c *fiber.Ctx) error {
	req := new(dto.GetInstitutionsReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	institutions, err := p.InstitutionRepository.GetAllInstitutions(req.EventId)
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

func (p *ParticipantDataController) AddInstitution(c *fiber.Ctx) error {
	req := new(dto.AddInstitutionReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	institution := &entity.Institution{
		EventId:         req.EventId,
		Name:            req.Name,
		Email:           req.Email,
		PendampingName:  req.PendampingName,
		PendampingPhone: req.PendampingPhone,
		Issuer:          entity.USER,
	}

	institution, err := p.InstitutionRepository.CreateInstitution(institution)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding institution"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   institution,
	})
}

func (p *ParticipantDataController) PickInstitution(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)
	req := new(dto.PickInstitutionReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	currentParticipant, err := p.ParticipantRepository.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting participant"),
		})
	}

	valid := currentParticipant.ValidateUserStep(entity.StepSelectInstitutionParticipant)
	if !valid {
		return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Anda tidak dapat lagi update data ini"),
		})
	}

	participant := entity.Participant{
		InstitutionId: req.InstitutionId,
		ProgressStep:  entity.StepSelectInstitutionParticipant,
	}

	if err := p.ParticipantRepository.UpdateParticipant(participantId, &participant); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating participant"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}

func (p *ParticipantDataController) AddMembers(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)
	req := new(dto.AddMemberReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	currentParticipant, err := p.ParticipantRepository.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting participant"),
		})
	}

	valid := currentParticipant.ValidateUserStep(entity.StepFillBiodatasParticipant)
	if !valid {
		return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Anda tidak dapat lagi update data ini"),
		})
	}

	biodataCreate := &entity.Biodata{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		IdNumber: req.IDNumber,
	}

	file, _ := c.FormFile("id_card")
	if file != nil {
		if err := util.ValidateFile(file); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		idCardUrl, err := util.FileSaver(file, "id_card"+participantId, "id_card/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		biodataCreate.IdCardPicture = idCardUrl
	}

	biodata, err := p.ParticipantRepository.AddBiodata(participantId, biodataCreate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding member"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   biodata,
	})
}

// TODO: Remove Members
func (p *ParticipantDataController) RemoveMembers(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)
	biodataId := c.Params("id")

	currentParticipant, err := p.ParticipantRepository.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting participant"),
		})
	}

	valid := currentParticipant.ValidateUserStep(entity.StepFillBiodatasParticipant)
	if !valid {
		return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Anda tidak dapat lagi update data ini"),
		})
	}

	if err := p.ParticipantRepository.RemoveBiodata(biodataId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when removing member"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}

func (p *ParticipantDataController) LockData(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	currentParticipant, err := p.ParticipantRepository.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting participant"),
		})
	}

	valid := currentParticipant.ValidateUserStep(entity.StepLockedParticipant)
	if !valid {
		return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Anda tidak dapat lagi update data ini"),
		})
	}

	cTime := time.Now()
	participant := entity.Participant{
		LockedAt:     &cTime,
		ProgressStep: entity.StepLockedParticipant,
	}

	if err := p.ParticipantRepository.UpdateParticipant(participantId, &participant); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating participant"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}
