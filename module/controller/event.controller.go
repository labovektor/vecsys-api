package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/module/dto"
	repository "github.com/labovector/vecsys-api/module/repository/event"
	"github.com/labovector/vecsys-api/util"
)

type EventController struct {
	eventRepo repository.EventRepository
}

func NewEventController(eventRepo repository.EventRepository) EventController {
	return EventController{
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

	id, _ := util.GetIdSession(c)

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
	event, err := ec.eventRepo.FindAllEvent()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Error{
			Message: "Something wrong when getting event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(event)
}
