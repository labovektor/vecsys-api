package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	ir "github.com/labovector/vecsys-api/internal/rest/repository/institution"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
	"gorm.io/gorm"
)

type ParticipantDataController struct {
	ParticipantRepository ur.UserRepository
	InstitutionRepository ir.InstitutionRepository

	// In case we need custom tx
	db *gorm.DB
}

func NewParticipantDataController(
	participantRepository ur.UserRepository,
	institutionRepository ir.InstitutionRepository,
	db *gorm.DB,
) *ParticipantDataController {
	return &ParticipantDataController{
		ParticipantRepository: participantRepository,
		InstitutionRepository: institutionRepository,
		db:                    db,
	}
}

func (p *ParticipantDataController) GetAllInstitution(c *fiber.Ctx) error {
	eventId := c.Locals(util.CurentUserEventIdKey).(string)

	institutions, err := p.InstitutionRepository.GetAllInstitutions(eventId)
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

	participantId := c.Locals(util.CurrentUserIdKey).(string)
	eventId := c.Locals(util.CurentUserEventIdKey).(string)

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

	institution := &entity.Institution{
		EventId:         &eventId,
		Name:            req.Name,
		Email:           req.Email,
		PendampingName:  req.PendampingName,
		PendampingPhone: req.PendampingPhone,
		Issuer:          entity.USER,
	}

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	institution, err := p.InstitutionRepository.WithDB(tx).CreateInstitution(institution)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding institution"),
		})
	}

	institutionId := institution.Id.String()

	participant := entity.Participant{
		InstitutionId: &institutionId,
		ProgressStep:  entity.StepSelectInstitutionParticipant,
	}

	if err := p.ParticipantRepository.WithDB(tx).UpdateParticipant(participantId, &participant); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating participant"),
		})
	}

	tx.Commit()
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
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	participant := entity.Participant{
		InstitutionId: &req.InstitutionId,
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
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
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

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	biodata, err := p.ParticipantRepository.WithDB(tx).AddBiodata(&participantId, biodataCreate)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding member"),
		})
	}

	participant := entity.Participant{
		ProgressStep: entity.StepFillBiodatasParticipant,
	}

	if err := p.ParticipantRepository.WithDB(tx).UpdateParticipant(participantId, &participant); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating participant"),
		})
	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   biodata,
	})
}

// TODO: Remove Members
func (p *ParticipantDataController) RemoveMembers(c *fiber.Ctx) error {
	biodataId := c.Params("id")

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

func (p *ParticipantDataController) GenerateCard(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	participant, err := p.ParticipantRepository.FindParticipantById(participantId, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	if !participant.IsLocked() {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("User belum mengunci semua data!"),
		})
	}

	cardBytes, err := util.GenerateCard(participant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat membuat kartu peserta"),
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"kartu_peserta_%s.pdf\"", participant.Name))
	c.Set("Content-Length", fmt.Sprintf("%d", len(cardBytes)))

	return c.Send(cardBytes)
}
