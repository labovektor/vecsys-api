package controller

import (
	"github.com/gofiber/fiber/v2"
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
	emailSession, _ := util.GetEmailSession(c)

	participant, err := ac.userRepo.FindParticipantByEmail(emailSession)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat mengambil data user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(participant)
}
