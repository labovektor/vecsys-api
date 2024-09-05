package controller

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/module/admin/repository"
	"github.com/labovector/vecsys-api/module/admin/usecase"
	"github.com/labovector/vecsys-api/util"
)

type authController struct {
	adminRepo repository.AdminRepository
}

func NewAuthController(adminRepo repository.AdminRepository) *authController {
	return &authController{
		adminRepo: adminRepo,
	}
}

func (ac *authController) Login(c *fiber.Ctx) error {
	req := new(entity.LoginReq)

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

	if err := usecase.GenerateSession(c, admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat membuat sesi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(admin)

}

func (ac *authController) Register(c *fiber.Ctx) error {
	req := new(entity.SignUpReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
			Message: "can't handle request",
		})
	}

	_, err := ac.adminRepo.FindAdminByUsername(req.Username)
	if err == nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
			Message: "Username tersebut telah dipakai",
		})
	}

	emailValid := util.ValidateEmail(req.Email)
	if !emailValid {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
			Message: "Email tidak valid",
		})
	}

	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: err.Error(),
		})
	}

	// Mengambil file dari form-data
	file, _ := c.FormFile("profile_picture")

	var newFileName string

	if file != nil {
		// Generate nama file baru
		newFileName = ac.generateNewFilename(req.Username, file.Filename)

		filePath := filepath.Join("./__public/profile/", newFileName)

		// Membuat direktori jika belum ada
		err := os.MkdirAll("./__public/profile/", os.ModePerm)
		if err != nil {
			return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
				Message: "Error when processing data",
			})
		}

		// Simpan file ke server
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Error{
				Message: "Error when processing data",
			})
		}
	}

	admin := entity.Admin{
		Username:       req.Username,
		DisplayName:    req.DisplayName,
		ProfilePicture: filepath.Join("/public", newFileName),
		Email:          req.Email,
		Password:       passwordHash,
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

func (ac *authController) Logout(c *fiber.Ctx) error {
	if err := usecase.InvalidateSession(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Kesalahan saat menghapus sesi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})
}

func (ac *authController) generateNewFilename(userID string, filename string) string {
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s%s", userID, ext)
	return newFilename
}
