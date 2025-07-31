package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	r "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
)

func UserStepMiddleware(repository r.UserRepository, targetStep entity.ParticipantProgress) fiber.Handler {
	return func(c *fiber.Ctx) error {
		participantId := c.Locals(util.CurrentUserIdKey).(string)

		currentParticipant, err := repository.FindParticipantById(participantId, false)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Something wrong when getting participant"),
			})
		}

		valid := currentParticipant.ValidateUserStep(targetStep)
		if !valid {
			return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Anda tidak dapat lagi update data ini"),
			})
		}

		return c.Next()
	}
}
