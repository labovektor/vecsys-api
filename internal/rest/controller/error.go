package controller

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repoerrors "github.com/labovector/vecsys-api/internal/rest/repository"
)

// respondError maps ErrNotFound to 404, logs every unexpected error, and keeps
// internal details out of the response body.
func respondError(c *fiber.Ctx, err error, notFoundMessage, internalMessage string) error {
	if errors.Is(err, repoerrors.ErrNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(notFoundMessage),
		})
	}

	log.Printf("%s %s: %s: %v", c.Method(), c.Path(), internalMessage, err)
	return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
		Status: dto.ErrorStatus.WithMessage(internalMessage),
	})
}
