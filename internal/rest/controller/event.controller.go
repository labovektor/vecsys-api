package controller

import (
	"path/filepath"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
			Message: "Masukkan data dengan benar",
		})
	}

	id := c.Locals(util.CurrentUserIdKey).(string)

	event := entity.Event{
		AdminId: id,
		Name:    req.Name,
	}

	event, err := ec.eventRepo.CreateEvent(&event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when creating event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(event)
}

func (ec *EventController) GetAllEvent(c *fiber.Ctx) error {
	id := c.Locals(util.CurrentUserIdKey).(string)
	event, err := ec.eventRepo.FindAllEvent(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when getting event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(event)
}

func (ec *EventController) GetEventById(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := ec.eventRepo.FindEventById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when getting event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(event)
}

func (ec *EventController) DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	err := ec.eventRepo.DeleteEvent(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when deleting event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})
}

func (ec *EventController) UpdateEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(dto.EventEditReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
			Message: "Masukkan data dengan benar",
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

		iconUrl, err := util.FileSaver(file, id, "event/")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
				Message: err.Error(),
			})
		}

		event.Icon = iconUrl
	}

	err := ec.eventRepo.UpdateEvent(id, &event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when updating event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})
}

func (ec *EventController) ToggleEventActive(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := ec.eventRepo.FindEventById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when getting event",
		})
	}

	event = &entity.Event{
		Active: event.Active,
	}

	err = ec.eventRepo.UpdateEvent(id, event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when updating event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status": "Success",
	})
}
