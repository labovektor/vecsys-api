package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/internal/util"
)

func UserMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := util.ValidateSessionUser(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
				Message: err.Error(),
			})
		}

		return c.Next()
	}
}
