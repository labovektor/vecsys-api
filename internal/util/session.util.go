package util

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labovector/vecsys-api/entity"
)

const (
	ROLE_ADMIN       = "role_admin"
	ROLE_USER        = "role_user"
	EXP              = 24 * time.Hour
	CurrentUserIdKey = "current_user_id"

	// in case of user the email will be used for current_user_name value
	CurrentUserNameKey = "current_user_name"

	// This will only available for user (participant), not for admin
	CurentUserEventIdKey = "current_user_event_id"
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
	sess.Set("event_id", *participant.EventId)

	// Save session
	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to create session")
	}

	return nil
}

func ValidateSessionAdmin(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

	username := sess.Get("username")
	id := sess.Get("id")
	role := sess.Get("role")
	if username == nil || id == nil || role == nil {
		return fmt.Errorf("failed to get session")
	}

	if role.(string) != ROLE_ADMIN {
		return fmt.Errorf("role is not admin")
	}

	sess.SetExpiry(EXP)

	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to save session")
	}

	c.Locals(CurrentUserIdKey, id)
	c.Locals(CurrentUserNameKey, username)

	return nil
}

func ValidateSessionUser(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

	username := sess.Get("email")
	id := sess.Get("id")
	eventId := sess.Get("event_id")
	role := sess.Get("role")
	if username == nil || id == nil || role == nil {
		return fmt.Errorf("failed to get session")
	}
	sess.SetExpiry(EXP)

	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to save session")
	}

	c.Locals(CurrentUserIdKey, id)
	c.Locals(CurrentUserNameKey, username)
	c.Locals(CurentUserEventIdKey, eventId)

	return nil
}

func InvalidateSession(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

	if err := sess.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy session")
	}

	return nil
}
