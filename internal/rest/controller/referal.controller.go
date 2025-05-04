package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/referal"
	"github.com/labovector/vecsys-api/internal/util"
)

type ReferalController struct {
	referalRepo repository.ReferalRepository
}

func NewReferalController(referalRepo repository.ReferalRepository) *ReferalController {
	return &ReferalController{
		referalRepo: referalRepo,
	}
}

func (r *ReferalController) GetReferalsByEventId(c *fiber.Ctx) error {
	eventId := c.Params("id")
	referals, err := r.referalRepo.GetAllVoucher(eventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting referals"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   referals,
	})
}

func (r *ReferalController) GetReferalByCode(c *fiber.Ctx) error {
	referalCode := c.Params("code")
	referal, err := r.referalRepo.GetVoucherByCode(referalCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting referal"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   referal,
	})
}

func (r *ReferalController) AddReferalToEvent(c *fiber.Ctx) error {
	eventId := c.Params("id")
	req := new(dto.ReferalCreateReq)
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

	referal := entity.Referal{
		EventId:       &eventId,
		Code:          req.Code,
		Desc:          req.Desc,
		SeatAvailable: *req.SeatAvailable,
		IsDiscount:    req.IsDiscount,
		Discount:      req.Discount,
	}

	referal, err := r.referalRepo.CreateVoucher(&referal)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding referal"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   referal,
	})
}

func (r *ReferalController) DeleteReferal(c *fiber.Ctx) error {
	referalId := c.Params("id")
	err := r.referalRepo.DeleteVoucher(referalId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when deleting referal"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}
