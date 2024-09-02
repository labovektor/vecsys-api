package usecase

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labovector/vecsys-api/entity"
)

const (
	ROLE_USER = "role_user"
)

func GenerateSession(c *fiber.Ctx, admin *entity.Admin) error {
	sess := c.Locals("session").(*session.Session)

	sess.Set("username", admin.Username)
	sess.Set("role", ROLE_USER)

	// Save session
	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to create session")
	}

	return nil
}

func RegenerateSession(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

	// Extend the session for the next 1 hour
	sess.SetExpiry(1 * time.Hour)

	// Regenerate session
	if err := sess.Regenerate(); err != nil {
		return fmt.Errorf("failed to create session")
	}

	// Save session
	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to create session")
	}

	return nil
}

func ValidateSession(c *fiber.Ctx) bool {
	sess := c.Locals("session").(*session.Session)

	username := sess.Get("username")
	return username != nil
}

func GetUsernameSession(c *fiber.Ctx) string {
	sess := c.Locals("session").(*session.Session)

	username := sess.Get("username")
	if username == nil {
		return ""
	}

	return username.(string)
}

func InvalidateSession(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

	if err := sess.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy session")
	}

	return nil
}
