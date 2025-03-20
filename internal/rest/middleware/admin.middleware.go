package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/internal/util"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		valid := util.ValidateSessionAdmin(c)
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
				Message: "You are not allowed to access this",
			})
		}

		err := c.Next()
		if err != nil {
			return err
		}

		if err := util.RegenerateSession(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
				Message: "Kesalahan saat membuat sesi",
			})
		}

		return nil
	}
}

// Strict session means you cannot extend your session time (only valid for 1 hour long)
func AdminMiddlewareStrictSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		valid := util.ValidateSessionAdmin(c)
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
				Message: "You are not allowed to access this",
			})
		}

		return c.Next()
	}
}
