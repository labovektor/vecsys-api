package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/internal/util"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := util.ValidateSessionAdmin(c); err != nil {
			if errors.Is(err, util.ErrUnauthorized) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
					Code:    fiber.StatusUnauthorized,
					Message: "Sesi tidak valid",
				})
			}

			log.Printf("%s %s: session validation: %v", c.Method(), c.Path(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: "Kesalahan internal",
			})
		}
		return c.Next()
	}
}
