package controller

import (
	"github.com/gofiber/fiber/v2"
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
