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

	// In case we need custom tx
	db *gorm.DB
}

func NewParticipantController(
	participantRepository ur.UserRepository,
	categoryRepository cr.CategoryRepository,
	regionRepository rr.RegionRepository,
	paymentRepository pr.PaymentRepository,
	referalRepository vr.ReferalRepository,
	db *gorm.DB,
) *ParticipantAdministrationController {
	return &ParticipantAdministrationController{
		ParticipantRepository: participantRepository,
		CategoryRepository:    categoryRepository,
		RegionRepository:      regionRepository,
		PaymentRepository:     paymentRepository,
		ReferalRepository:     referalRepository,
		db:                    db,
	}
}

func (p *ParticipantAdministrationController) GetAllEventCategoryAndRegion(c *fiber.Ctx) error {
	eventId := c.Locals(util.CurentUserEventIdKey).(string)

	categories, err := p.CategoryRepository.GetAllActiveCategories(eventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting categories"),
		})
	}

	regions, err := p.RegionRepository.GetAllActiveRegion(eventId)
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
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
		})
	}

	if err := util.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	participant := entity.Participant{
		CategoryId:   &req.CategoryId,
		RegionId:     &req.RegionId,
		ProgressStep: entity.StepCategorizedParticipant,
	}

	if err := p.ParticipantRepository.UpdateParticipant(participantId, &participant); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating participant"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}

func (p *ParticipantAdministrationController) GetAllPaymentOption(c *fiber.Ctx) error {
	eventId := c.Locals(util.CurentUserEventIdKey).(string)

	paymentOptions, err := p.PaymentRepository.GetPaymentOptions(eventId)
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
	eventId := c.Locals(util.CurentUserEventIdKey).(string)
	req := new(dto.ClaimReferalReq)

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

	referal, err := p.ReferalRepository.GetVoucherByCode(req.Code, eventId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Referal tidak valid"),
		})
	}

	if !referal.IsAvailableToClaim() {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Referal sudah tidak tersedia"),
		})
	}

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah partisipan sudah puya paymen
	// Jika belum punya maka buatkan payment baru
	// Jika sudah punya maka cek apakah paymentnya punya referal
	currentParticipantPayment, err := p.PaymentRepository.WithDB(tx).GetPaymentByParticipantId(participantId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			referalId := referal.Id.String()
			payment := entity.Payment{
				ParticipantId: &participantId,
				ReferalId:     &referalId,
			}

			_, err = p.PaymentRepository.WithDB(tx).CreatePayment(&payment)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
					Status: dto.ErrorStatus.WithMessage("Something wrong when creating payment"),
				})
			}

			updatedRef := entity.Referal{
				SeatAvailable: referal.SeatAvailable - 1,
			}
			if err = p.ReferalRepository.WithDB(tx).UpdateVoucher(referal.Id.String(), &updatedRef); err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
					Status: dto.ErrorStatus.WithMessage("Something wrong when updating referal"),
				})
			}

			tx.Commit()
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

	if currentParticipantPayment.ReferalId != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.SuccessStatus.WithMessage("Anda sudah memiliki referal"),
			Data:   referal,
		})
	}

	referalId := referal.Id.String()
	payment := entity.Payment{
		ReferalId: &referalId,
	}

	if err := p.PaymentRepository.WithDB(tx).UpdatePayment(currentParticipantPayment.Id.String(), &payment); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating payment"),
		})
	}

	updatedRef := entity.Referal{
		SeatAvailable: referal.SeatAvailable - 1,
	}
	if err = p.ReferalRepository.WithDB(tx).UpdateVoucher(referal.Id.String(), &updatedRef); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating referal"),
		})
	}

	tx.Commit()
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
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar!"),
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
		ParticipantId:   &participantId,
		PaymentOptionId: &req.PaymentOptionId,
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

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah partisipan sudah puya paymen
	// Jika belum punya maka buatkan payment baru
	// Jika sudah punya maka update paymentnya
	currentParticipantPayment, err := p.PaymentRepository.WithDB(tx).GetPaymentByParticipantId(participantId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, err = p.PaymentRepository.WithDB(tx).CreatePayment(&payment)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
					Status: dto.ErrorStatus.WithMessage("Something wrong when creating payment"),
				})
			}
		} else {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Something wrong when getting payment"),
			})
		}
	}

	if err = p.PaymentRepository.WithDB(tx).UpdatePayment(currentParticipantPayment.Id.String(), &payment); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating payment"),
		})
	}

	participant := entity.Participant{
		ProgressStep: entity.StepPaidParticipant,
	}

	if err := p.ParticipantRepository.WithDB(tx).UpdateParticipant(participantId, &participant); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating participant"),
		})
	}

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}
