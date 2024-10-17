package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/repository"
	"github.com/labovector/vecsys-api/util"
)

type userController struct {
	userRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) userController {
	return userController{
		userRepo: userRepo,
	}
}

func (ac *userController) GetUser(c *fiber.Ctx) error {
	emailSession, err := util.GetEmailSession(c)
	if err == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
			Message: "Sesi tidak valid",
		})
	}

	participant, err := ac.userRepo.FindParticipantByEmail(emailSession)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat mengambil data user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(participant)
}
