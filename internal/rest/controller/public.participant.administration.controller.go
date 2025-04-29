package controller

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	cr "github.com/labovector/vecsys-api/internal/rest/repository/category"
	pr "github.com/labovector/vecsys-api/internal/rest/repository/payment"
	vr "github.com/labovector/vecsys-api/internal/rest/repository/referal"
	rr "github.com/labovector/vecsys-api/internal/rest/repository/region"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
	"gorm.io/gorm"
)

type ParticipantAdministrationController struct {
	ParticipantRepository ur.UserRepository
	CategoryRepository    cr.CategoryRepository
	RegionRepository      rr.RegionRepository
	PaymentRepository     pr.PaymentRepository
	ReferalRepository     vr.ReferalRepository
}

func NewParticipantController(
	participantRepository ur.UserRepository,
	categoryRepository cr.CategoryRepository,
	regionRepository rr.RegionRepository,
	paymentRepository pr.PaymentRepository,
	referalRepository vr.ReferalRepository,
) *ParticipantAdministrationController {
	return &ParticipantAdministrationController{
		ParticipantRepository: participantRepository,
		CategoryRepository:    categoryRepository,
		RegionRepository:      regionRepository,
		PaymentRepository:     paymentRepository,
		ReferalRepository:     referalRepository,
	}
}

func (p *ParticipantAdministrationController) GetAllEventCategoryAndRegion(c *fiber.Ctx) error {
	req := new(dto.GetAllEventCategoryAndRegionReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}
	categories, err := p.CategoryRepository.GetAllCategories(req.EventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting categories"),
		})
	}

	regions, err := p.RegionRepository.GetAllRegion(req.EventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting regions"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data: fiber.Map{
			"categories": categories,
			"regions":    regions,
		},
	})
}

func (p *ParticipantAdministrationController) PickCategoryAndRegion(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	req := new(dto.PickCategoryAndRegionReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	currentParticipant, err := p.ParticipantRepository.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting participant"),
		})
	}

	valid := currentParticipant.ValidateUserStep(entity.StepCategorizedParticipant)
	if !valid {
		return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Anda tidak dapat lagi update data ini"),
		})
	}

	participant := entity.Participant{
		CategoryId:   req.CategoryId,
		RegionId:     req.RegionId,
		ProgressStep: entity.StepCategorizedParticipant,
	}

	if err = p.ParticipantRepository.UpdateParticipant(participantId, &participant); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating participant"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}

func (p *ParticipantAdministrationController) GetAllPaymentOption(c *fiber.Ctx) error {
	req := new(dto.GetPaymentOptionsReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	paymentOptions, err := p.PaymentRepository.GetPaymentOptions(req.EventId)
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

// TODO: Validate Referal
func (p *ParticipantAdministrationController) ValidateReferal(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)
	req := new(dto.ClaimReferalReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	referal, err := p.ReferalRepository.GetVoucherByCode(req.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Referal tidak valid"),
		})
	}

	available := referal.IsAvailableToClaim()
	if !available {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Referal sudah tidak tersedia"),
		})
	}

	// Cek apakah partisipan sudah puya paymen
	// Jika belum punya maka buatkan payment baru
	// Jika sudah punya maka cek apakah paymentnya punya referal
	currentParticipantPayment, err := p.PaymentRepository.GetPaymentByParticipantId(participantId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			payment := entity.Payment{
				ParticipantId: participantId,
				ReferalId:     referal.Id.String(),
			}

			_, err = p.PaymentRepository.CreatePayment(&payment)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
					Status: dto.ErrorStatus.WithMessage("Something wrong when creating payment"),
				})
			}

			updatedRef := entity.Referal{
				SeatAvailable: referal.SeatAvailable - 1,
			}
			if err = p.ReferalRepository.UpdateVoucher(referal.Id.String(), &updatedRef); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
					Status: dto.ErrorStatus.WithMessage("Something wrong when updating referal"),
				})
			}

			return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
				Status: dto.SuccessStatus,
				Data:   referal,
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Something wrong when getting payment"),
			})
		}
	}

	if currentParticipantPayment.ReferalId != "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Anda sudah memiliki referal"),
		})
	}

	updatedRef := entity.Referal{
		SeatAvailable: referal.SeatAvailable - 1,
	}
	if err = p.ReferalRepository.UpdateVoucher(referal.Id.String(), &updatedRef); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating referal"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   referal,
	})
}

// TODO: Payment
func (p *ParticipantAdministrationController) Payment(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)
	req := new(dto.PaymentReq)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Gagal Memproses Data!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	transferDate, err := time.Parse("2006-01-02T15:04:05Z07:00", req.TransferDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}
	payment := entity.Payment{
		ParticipantId:   participantId,
		PaymentOptionId: req.PaymentOptionId,
		BankName:        req.AccountName,
		BankAccount:     req.AccountNumber,
		Date:            &transferDate,
	}

	file, _ := c.FormFile("invoice")
	if file != nil {
		if err := util.ValidateFile(file); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		invoiceUrl, err := util.FileSaver(file, "invoice"+participantId, "invoice/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		payment.Invoice = invoiceUrl
	}

	// Cek apakah partisipan sudah puya paymen
	// Jika belum punya maka buatkan payment baru
	// Jika sudah punya maka update paymentnya
	currentParticipantPayment, err := p.PaymentRepository.GetPaymentByParticipantId(participantId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, err = p.PaymentRepository.CreatePayment(&payment)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
					Status: dto.ErrorStatus.WithMessage("Something wrong when creating payment"),
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Something wrong when getting payment"),
			})
		}
	}

	if err = p.PaymentRepository.UpdatePayment(currentParticipantPayment.Id.String(), &payment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating payment"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}
