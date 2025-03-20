package util

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labovector/vecsys-api/entity"
)

const (
	ROLE_ADMIN = "role_admin"
	ROLE_USER  = "role_user"
)

func GenerateSessionAdmin(c *fiber.Ctx, admin *entity.Admin) error {
	sess := c.Locals("session").(*session.Session)

	sess.Set("username", admin.Username)
	sess.Set("id", admin.Id.String())
	sess.Set("role", ROLE_ADMIN)

	// Save session
	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to create session")
	}

	return nil
}

func GenerateSessionUser(c *fiber.Ctx, participant *entity.Participant) error {
	sess := c.Locals("session").(*session.Session)

	sess.Set("email", participant.Email)
	sess.Set("id", participant.Id.String())
	sess.Set("role", ROLE_USER)

	// Save session
	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to create session")
	}

	return nil
}

func RegenerateSession(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

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

func ValidateSessionAdmin(c *fiber.Ctx) bool {
	sess := c.Locals("session").(*session.Session)

	username := sess.Get("username")
	id := sess.Get("id")
	role := sess.Get("role")
	if username == nil || id == nil || role == nil {
		return false
	}

	return role.(string) == ROLE_ADMIN
}

func ValidateSessionUser(c *fiber.Ctx) bool {
	sess := c.Locals("session").(*session.Session)

	username := sess.Get("email")
	id := sess.Get("id")
	role := sess.Get("role")
	return username != nil && id != nil && role != nil
}

func GetEmailSession(c *fiber.Ctx) (string, error) {
	sess := c.Locals("session").(*session.Session)

	email := sess.Get("email")
	if email == nil {
		return "", fmt.Errorf("failed to get email from session")
	}
	return email.(string), nil
}

func GetUsernameSession(c *fiber.Ctx) (string, error) {
	sess := c.Locals("session").(*session.Session)

	username := sess.Get("username")
	if username == nil {
		return "", fmt.Errorf("failed to get username from session")
	}
	return username.(string), nil
}

func GetIdSession(c *fiber.Ctx) (string, error) {
	sess := c.Locals("session").(*session.Session)

	id := sess.Get("id")
	if id == nil {
		return "", &fiber.Error{Message: "Username Not Found"}
	}

	return id.(string), nil
}

func InvalidateSession(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

	if err := sess.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy session")
	}

	return nil
}
