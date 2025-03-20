package controller

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	adminRepo "github.com/labovector/vecsys-api/internal/rest/repository/admin"
	userRepo "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
)

type AuthController struct {
	adminRepo adminRepo.AdminRepository
	userRepo  userRepo.UserRepository
}

func NewAuthController(adminRepo adminRepo.AdminRepository, userRepo userRepo.UserRepository) *AuthController {
	return &AuthController{
		adminRepo: adminRepo,
		userRepo:  userRepo,
	}
}

func (ac *AuthController) LoginAdmin(c *fiber.Ctx) error {
	req := new(dto.AdminLoginReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
			Message: "can't handle request",
		})
	}

	admin, err := ac.adminRepo.FindAdminByUsername(req.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Error{
			Message: "Username tidak ditemukan",
		})
	}

	match := util.CheckPasswordHash(req.Password, admin.Password)
	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
			Message: "Password salah",
		})
	}

	if err := util.GenerateSessionAdmin(c, admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat membuat sesi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(admin)

}

func (ac *AuthController) RegisterAdmin(c *fiber.Ctx) error {
	req := new(dto.AdminSignUpReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
			Message: "can't handle request",
		})
	}

	_, err := ac.adminRepo.FindAdminByUsername(req.Username)
	if err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Username tersebut telah dipakai",
		})
	}

	emailValid := util.ValidateEmail(req.Email)
	if !emailValid {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Email tidak valid",
		})
	}

	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: err.Error(),
		})
	}

	admin := entity.Admin{
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Password:    passwordHash,
	}

	file, _ := c.FormFile("profile_picture")

	if file != nil {
		if file.Size > 10*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
				Message: "Max file size 10MB",
			})
		}
		ext := filepath.Ext(file.Filename)
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
				Message: "Only accept image file",
			})
		}

		profileUrl, err := util.FileSaver(file, admin.Username, "profile/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
				Message: err.Error(),
			})
		}

		admin.ProfilePicture = profileUrl
	}

	_, err = ac.adminRepo.CreateAdmin(&admin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when creating user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})

}

func (ac *AuthController) LogoutAdmin(c *fiber.Ctx) error {
	if err := util.InvalidateSession(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat menghapus sesi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})
}

func (ac *AuthController) LoginUser(c *fiber.Ctx) error {
	req := new(dto.ParticipantLoginReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
			Message: "can't handle request",
		})
	}

	user, err := ac.userRepo.FindParticipantByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Error{
			Message: "Username tidak ditemukan",
		})
	}

	match := util.CheckPasswordHash(req.Password, user.Password)
	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
			Message: "Password salah",
		})
	}

	if err := util.GenerateSessionUser(c, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat membuat sesi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)

}

func (ac *AuthController) RegisterUser(c *fiber.Ctx) error {
	req := new(dto.ParticipantSignUpReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
			Message: "can't handle request",
		})
	}

	_, err := ac.userRepo.FindParticipantByEmail(req.Email)
	if err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Email tersebut telah dipakai",
		})
	}

	emailValid := util.ValidateEmail(req.Email)
	if !emailValid {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Email tidak valid",
		})
	}

	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: err.Error(),
		})
	}

	participant := entity.Participant{
		EventId:  req.EventId,
		Email:    req.Email,
		Name:     req.Name,
		Password: passwordHash,
	}

	_, err = ac.userRepo.CreateParticipant(&participant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when creating user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})

}

func (ac *AuthController) LogoutUser(c *fiber.Ctx) error {
	if err := util.InvalidateSession(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat menghapus sesi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})
}
