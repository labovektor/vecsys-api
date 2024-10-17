package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/util"
)

func UserMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		valid := util.ValidateSessionUser(c)
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
				Message: "You are not allowed to access this",
			})
		}

		// if err := util.RegenerateSession(c); err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
		// 		Message: "Kesalahan saat membuat sesi",
		// 	})
		// }

		return c.Next()
	}
}

// Strict session means you cannot extend your session time (only valid for 1 hour long)
func UserMiddlewareStrictSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		valid := util.ValidateSessionUser(c)
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
				Message: "You are not allowed to access this",
			})
		}

		return c.Next()
	}
}
