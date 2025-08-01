package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/internal/rest/dto"
	repository "github.com/labovector/vecsys-api/internal/rest/repository/event"
	"github.com/labovector/vecsys-api/internal/util"
)

type EventController struct {
	eventRepo repository.EventRepository
}

func NewEventController(eventRepo repository.EventRepository) *EventController {
	return &EventController{
		eventRepo: eventRepo,
	}
}

func (ec *EventController) CreateEvent(c *fiber.Ctx) error {
	req := new(dto.EventCreateReq)
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

	id := c.Locals(util.CurrentUserIdKey).(string)

	event := entity.Event{
		AdminId: &id,
		Name:    req.Name,
	}

	event, err := ec.eventRepo.CreateEvent(&event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when creating event"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   event,
	})
}

func (ec *EventController) GetAllEvent(c *fiber.Ctx) error {
	id := c.Locals(util.CurrentUserIdKey).(string)
	event, err := ec.eventRepo.FindAllEvent(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting event"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   event,
	})
}

func (ec *EventController) GetEventById(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := ec.eventRepo.FindEventById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting event"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
		Data:   event,
	})
}

func (ec *EventController) DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	err := ec.eventRepo.DeleteEvent(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when deleting event"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}

func (ec *EventController) UpdateEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(dto.EventEditReq)
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

	event := entity.Event{
		Name:              req.Name,
		Desc:              req.Desc,
		GroupMemberNum:    req.GroupMemberNum,
		ParticipantTarget: req.ParticipantTarget,
		Period:            req.Period,
	}

	file, _ := c.FormFile("icon")
	if file != nil {
		if err := util.ValidateFile(file); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		iconUrl, err := util.FileSaver(file, id, "event/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Status: dto.ErrorStatus.WithMessage(err.Error()),
			})
		}

		event.Icon = iconUrl
	}

	err := ec.eventRepo.UpdateEvent(id, &event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating event"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}

func (ec *EventController) ToggleEventActive(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := ec.eventRepo.FindEventById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when getting event"),
		})
	}

	newStatus := !*event.Active

	event = &entity.Event{
		Active: &newStatus,
	}

	err = ec.eventRepo.UpdateEvent(id, event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status: dto.ErrorStatus.WithMessage("Something wrong when updating event"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status: dto.SuccessStatus,
	})
}
