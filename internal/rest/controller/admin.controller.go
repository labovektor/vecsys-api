package controller

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/admin"
	"github.com/labovector/vecsys-api/internal/util"
)

type AdminController struct {
	adminRepo repository.AdminRepository
}

func NewAdminController(adminRepo repository.AdminRepository) *AdminController {
	return &AdminController{
		adminRepo: adminRepo,
	}
}

func (ac *AdminController) UpdateAdminProfile(c *fiber.Ctx) error {
	req := new(dto.AdminEditReq)

	id := c.Locals(util.CurrentUserIdKey).(string)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Failed to process data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	cAdmin, err := ac.adminRepo.FindAdminById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting user data"),
		})
	}

	admin := entity.Admin{
		DisplayName: req.DisplayName,
		Email:       req.Email,
	}

	file, _ := c.FormFile("profile_picture")

	if file != nil {
		if file.Size > 10*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Max file size is 10MB"),
			})
		}
		ext := filepath.Ext(file.Filename)
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Only accept image file"),
			})
		}

		profileUrl, err := util.FileSaver(file, cAdmin.Username, "profile/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		admin.ProfilePicture = profileUrl
	}

	err = ac.adminRepo.UpdateAdmin(id, &admin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating user data"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}

func (ac *AdminController) GetAdmin(c *fiber.Ctx) error {
	usernameSession := c.Locals(util.CurrentUserNameKey).(string)

	admin, err := ac.adminRepo.FindAdminByUsername(usernameSession)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting user data"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   admin,
	})
}
