package controller

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
)

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

func (ac *UserController) GetUser(c *fiber.Ctx) error {
	// Because in user, email is set as username
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	participant, err := ac.userRepo.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participant,
	})
}

func (ac *UserController) GetAllParticipantData(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	participant, err := ac.userRepo.FindParticipantById(participantId, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	res := fiber.Map{
		"participant": participant,
		"is_verified": participant.IsVerified(),
		"is_locked":   participant.IsLocked(),
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   res,
	})
}

func (ac *UserController) GetParticipantState(c *fiber.Ctx) error {
	participantId := c.Locals(util.CurrentUserIdKey).(string)

	participant, err := ac.userRepo.FindParticipantById(participantId, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	res := fiber.Map{
		"step":        participant.ProgressStep,
		"is_verified": participant.IsVerified(),
		"is_locked":   participant.IsLocked(),
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   res,
	})
}

func (ac *UserController) GetAllParticipant(c *fiber.Ctx) error {
	eventID := c.Params("id")
	step := c.Query("step", "all")

	var participants []entity.Participant
	var err error

	switch step {
	case "all":
		participants, err = ac.userRepo.FindAllParticipant(eventID)
	case "paid":
		participants, err = ac.userRepo.FindAllPaidParticipant(eventID)
	case "unpaid":
		participants, err = ac.userRepo.FindAllUnpaidParticipant(eventID)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Parameter step tidak valid"),
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participants,
	})
}

func (ac *UserController) GetParticipantByID(c *fiber.Ctx) error {
	ID := c.Params("id")
	participant, err := ac.userRepo.FindParticipantById(ID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participant,
	})
}

func (ac *UserController) VerifyParticipant(c *fiber.Ctx) error {
	ID := c.Params("id")
	participant, err := ac.userRepo.FindParticipantById(ID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	if participant.IsVerified() {
		participant.VerifiedAt = nil

		participant.ProgressStep = entity.StepPaidParticipant

		err = ac.userRepo.UpdateParticipant(ID, participant)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Kesalahan saat memperbarui data user"),
			})
		}
	} else {
		now := time.Now()
		participant.VerifiedAt = &now

		participant.ProgressStep = entity.StepVerifiedParticipant

		err = ac.userRepo.UpdateParticipant(ID, participant)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage("Kesalahan saat memperbarui data user"),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   participant,
	})
}

func (ac *UserController) GeneratePdfParticipant(c *fiber.Ctx) error {
	id := c.Params("id")
	participant, err := ac.userRepo.FindParticipantById(id, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengambil data user"),
		})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	iconPath := participant.Event.Icon
	iconPath = strings.ReplaceAll(iconPath, "\\", "/")
	iconURL := fmt.Sprintf("http://127.0.0.1:8787/api/v1%s", iconPath)
	cleanURL := strings.Split(iconURL, "?")[0]

	httpimg.Register(pdf, cleanURL, "")
	imgFormat := ""
	if strings.HasSuffix(strings.ToLower(cleanURL), ".jpg") {
		imgFormat = "JPG"
	} else if strings.HasSuffix(strings.ToLower(cleanURL), ".jpg") || strings.HasSuffix(strings.ToLower(cleanURL), ".jpeg") {
		imgFormat = "JPG"
	} else if strings.HasSuffix(strings.ToLower(cleanURL), ".gif") {
		imgFormat = "GIF"
	} else {
		imgFormat = "PNG"
	}
	pdf.Image(cleanURL, 90, 12, 30, 30, false, imgFormat, 0, "")
	pdf.Ln(50)
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 10, "Kartu Peserta", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, participant.Event.Name, "", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(50, 10, "Asal Sekolah", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Institution.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Jenjang/Kategori", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Category.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Region", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Region.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Nomor Peserta", "", 0, "", false, 0, "")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, ": "+participant.Id.String(), "", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(50, 10, "Nama", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Name, "", 1, "", false, 0, "")

	pdf.CellFormat(50, 10, "Email Peserta", "", 0, "", false, 0, "")
	pdf.CellFormat(0, 10, ": "+participant.Email, "", 1, "", false, 0, "")

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Anggota:")
	pdf.Ln(10)
	for i, m := range *participant.Biodata {
		pdf.SetFont("Arial", "", 14)
		pdf.CellFormat(0, 8, fmt.Sprintf("%d. %s (%s)", i+1, m.Name, m.IdNumber), "", 1, "", false, 0, "")
		pdf.Ln(5)
	}

	pdf.Ln(70)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Contact Person:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(0, 8, "089756620221"+" (Nami)", "", 1, "", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage(err.Error()),
		})
	}

	filename := fmt.Sprintf("kartu peserta_%s.pdf", participant.Name)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.SendStream(&buf)
}

func (ac *UserController) BulkAddParticipantFromCSV(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat mengupload file"),
		})
	}

	file, _ := fileHeader.Open()
	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat membaca file"),
		})
	}

	eventId := c.Params("id")
	if len(eventId) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Event ID kosong"),
		})
	}

	participants := make([]entity.Participant, 0, len(records)-1)
	for i, record := range records[1:] {
		if len(record) < 3 {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(fmt.Sprintf("Baris %d: Data tidak lengkap", i)),
			})
		}

		participantReq := dto.ParticipantSignUpReq{
			EventId:  eventId,
			Name:     record[0],
			Email:    record[1],
			Password: record[2],
		}

		if err := util.ValidateStruct(participantReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(fmt.Sprintf("Baris %d: %s", i, err.Error())),
			})
		}

		passwordHash, err := util.HashPassword(participantReq.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		participant := entity.Participant{
			EventId:  &participantReq.EventId,
			Name:     participantReq.Name,
			Email:    participantReq.Email,
			Password: passwordHash,
		}

		participants = append(participants, participant)
	}

	err = ac.userRepo.BulkAddParticipant(participants)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat menambahkan data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   fmt.Sprintf("%d data berhasil ditambahkan", len(participants)),
	})
}

func (ac *UserController) UpdateParticipantData(c *fiber.Ctx) error {
	participantId := c.Params("id")
	req := new(dto.ParticipantUpdateReq)
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

	participant := entity.Participant{
		Name:       *req.Name,
		CategoryId: req.CategoryId,
		RegionId:   req.RegionId,
	}

	err := ac.userRepo.UpdateParticipant(participantId, &participant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat memperbarui data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus.WithMessage("Data user berhasil diperbarui"),
		Data:   participant,
	})
}

func (ac *UserController) UpdateParticipantBiodata(c *fiber.Ctx) error {
	participantId := c.Params("id")

	req := new([]dto.BiodataUpdateReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Masukkan data dengan benar"),
		})
	}

	for _, biodata := range *req {
		if err := util.ValidateStruct(biodata); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}
	}

	biodatas := make([]entity.Biodata, 0, len(*req))
	for _, biodata := range *req {
		biodatas = append(biodatas, entity.Biodata{
			Id:       uuid.MustParse(biodata.Id),
			Name:     *biodata.Name,
			Gender:   *biodata.Gender,
			Email:    *biodata.Email,
			Phone:    *biodata.Phone,
			IdNumber: *biodata.IDNumber,
		})
	}

	err := ac.userRepo.BulkUpdateBiodata(participantId, biodatas)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat memperbarui data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus.WithMessage(fmt.Sprintf("%d data berhasil diperbarui", len(*req))),
	})

}

func (ac *UserController) DeleteParticipant(c *fiber.Ctx) error {
	participantId := c.Params("id")
	err := ac.userRepo.DeleteParticipant(participantId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Kesalahan saat menghapus data user"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus.WithMessage("Data user berhasil dihapus"),
	})
}
