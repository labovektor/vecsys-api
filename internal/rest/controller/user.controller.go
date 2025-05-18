package controller

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
)

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

func (ac *UserController) GetUser(c *fiber.Ctx) error {
	// Because in user, email is set as username
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	participant, err := ac.userRepo.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participant,
	})
}

func (ac *UserController) GetAllParticipantData(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	participant, err := ac.userRepo.FindParticipantById(participantId, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	res := fiber.Map{
		"participant": participant,
		"is_verified": participant.IsVerified(),
		"is_locked":   participant.IsLocked(),
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   res,
	})
}

func (ac *UserController) GetParticipantState(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	participant, err := ac.userRepo.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	res := fiber.Map{
		"step":        participant.ProgressStep,
		"is_verified": participant.IsVerified(),
		"is_locked":   participant.IsLocked(),
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   res,
	})
}

func (ac *UserController) GetAllParticipant(c *fiber.Ctx) error {
	eventId := c.Params("id")
	participant, err := ac.userRepo.FindAllParticipant(eventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participant,
	})
}

func (ac *UserController) GetParticipantByID(c *fiber.Ctx) error {
	ID := c.Params("id")
	participant, err := ac.userRepo.FindParticipantById(ID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participant,
	})
}

func (ac *UserController) VerifyParticipant(c *fiber.Ctx) error {
	ID := c.Params("id")
	participant, err := ac.userRepo.FindParticipantById(ID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	now := time.Now()
	participant.VerifiedAt = &now

	participant.ProgressStep = entity.StepVerifiedParticipant

	err = ac.userRepo.UpdateParticipant(ID, participant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat memperbarui data user"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participant,
	})
}

func (ac *UserController) BulkAddParticipantFromCSV(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengupload file"),
		})
	}

	file, _ := fileHeader.Open()
	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat membaca file"),
		})
	}

	eventId := c.Params("id")
	if len(eventId) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Event ID kosong"),
		})
	}

	participants := make([]entity.Participant, 0, len(records)-1)
	for i, record := range records[1:] {
		if len(record) < 3 {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(fmt.Sprintf("Baris %d: Data tidak lengkap", i)),
			})
		}

		participantReq := dto.ParticipantSignUpReq{
			EventId:  eventId,
			Name:     record[0],
			Email:    record[1],
			Password: record[2],
		}

		if err := util.ValidateStruct(participantReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(fmt.Sprintf("Baris %d: %s", i, err.Error())),
			})
		}

		participant := entity.Participant{
			EventId:  &participantReq.EventId,
			Name:     participantReq.Name,
			Email:    participantReq.Email,
			Password: participantReq.Password,
		}

		participants = append(participants, participant)
	}

	err = ac.userRepo.BulkAddParticipant(participants)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat menambahkan data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   fmt.Sprintf("%d data berhasil ditambahkan", len(participants)),
	})
}
