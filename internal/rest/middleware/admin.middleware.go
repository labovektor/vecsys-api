package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/internal/util"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := util.ValidateSessionAdmin(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
				Message: err.Error(),
			})
		}
		return c.Next()
	}
}
