package controller

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/module/dto"
	"github.com/labovector/vecsys-api/module/repository"
	"github.com/labovector/vecsys-api/util"
)

type adminController struct {
	adminRepo repository.AdminRepository
}

func NewAdminController(adminRepo repository.AdminRepository) *adminController {
	return &adminController{
		adminRepo: adminRepo,
	}
}

func (ac *adminController) UpdateAdminProfile(c *fiber.Ctx) error {
	req := new(dto.AdminEditReq)

	id, err := util.GetIdSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
			Message: err.Error(),
		})
	}

	cAdmin, err := ac.adminRepo.FindAdminById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat mengambil data user",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "can't handle request",
		})
	}

	admin := entity.Admin{
		DisplayName: req.DisplayName,
		Email:       req.Email,
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

		profileUrl, err := util.FileSaver(file, cAdmin.Username, "profile/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
				Message: err.Error(),
			})
		}

		admin.ProfilePicture = profileUrl
	}

	err = ac.adminRepo.UpdateAdmin(id, &admin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when updating user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})
}

func (ac *adminController) GetAdmin(c *fiber.Ctx) error {
	usernameSession, err := util.GetUsernameSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Error{
			Message: "Sesi tidak valid",
		})
	}

	admin, err := ac.adminRepo.FindAdminByUsername(usernameSession)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat mengambil data user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(admin)
}
