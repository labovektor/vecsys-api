package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/payment"
	"github.com/labovector/vecsys-api/internal/util"
)

type PaymentController struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentController(paymentRepo repository.PaymentRepository) PaymentController {
	return PaymentController{paymentRepo: paymentRepo}
}

func (cc *PaymentController) GetAllPaymentOptionByEventId(c *fiber.Ctx) error {
	eventId := c.Params("id")
	paymentOptions, err := cc.paymentRepo.GetPaymentOptions(eventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting payment options"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   paymentOptions,
	})
}

func (cc *PaymentController) GetPaymentOptionById(c *fiber.Ctx) error {
	paymentOptionId := c.Params("id")
	paymentOption, err := cc.paymentRepo.GetPaymentOptionById(paymentOptionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting payment option"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   paymentOption,
	})
}

func (cc *PaymentController) AddPaymentOptionToEvent(c *fiber.Ctx) error {
	eventId := c.Params("id")

	req := new(dto.PaymentOptionAddReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}
	payment := entity.PaymentOption{
		EventId:  &eventId,
		Provider: req.Provider,
		Account:  req.Account,
		Name:     req.Name,
		AsQR:     req.AsQR,
	}

	payment, err := cc.paymentRepo.CreatePaymentOption(&payment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when adding payment option"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   payment,
	})
}

func (cc *PaymentController) UpdatePaymentOption(c *fiber.Ctx) error {
	paymentOptionId := c.Params("id")
	req := new(dto.PaymentOptionUpdateReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
		})
	}

	payment := entity.PaymentOption{
		Provider: req.Provider,
		Account:  req.Account,
		Name:     req.Name,
		AsQR:     req.AsQR,
	}

	err := cc.paymentRepo.UpdatePaymentOption(paymentOptionId, &payment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating payment option"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   "Payment Updated Successfully",
	})

}

func (cc *PaymentController) DeletePaymentOption(c *fiber.Ctx) error {
	payemntOptionId := c.Params("id")
	err := cc.paymentRepo.DeletePaymentOption(payemntOptionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when deleting payment option"),
		})
	}
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}
