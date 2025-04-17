package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/category"
	"github.com/labovector/vecsys-api/internal/util"
)

type CategoryController struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryController(categoryRepo repository.CategoryRepository) CategoryController {
	return CategoryController{categoryRepo: categoryRepo}
}

func (cc *CategoryController) GetAllCategoryByEventId(c *fiber.Ctx) error {
	eventId := c.Params("id")
	categories, err := cc.categoryRepo.GetAllCategories(eventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting categories"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   categories,
	})
}

func (cc *CategoryController) GetCategoryById(c *fiber.Ctx) error {
	categoryId := c.Params("id")
	category, err := cc.categoryRepo.GetCategory(categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting category"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   category,
	})
}

func (cc *CategoryController) AddCategoryToEvent(c *fiber.Ctx) error {
	eventId := c.Params("id")
	req := new(dto.NewCategoryReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	category := entity.Category{
		EventId: eventId,
		Name:    req.Name,
		IsGroup: &req.IsGroup,
		Visible: &req.Visible,
	}

	category, err := cc.categoryRepo.CreateCategory(&category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when creating category"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   category,
	})
}

func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	categoryId := c.Params("id")
	req := new(dto.UpdateCategoryReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar"),
		})
	}

	category := entity.Category{
		Name:    req.Name,
		IsGroup: &req.IsGroup,
		Visible: &req.Visible,
	}

	category, err := cc.categoryRepo.UpdateCategory(categoryId, &category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating category"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   "Category Updated Successfully",
	})
}

func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	categoryId := c.Params("id")
	err := cc.categoryRepo.DeleteCategory(categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when deleting category"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}
