package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/controller"
	"github.com/labovector/vecsys-api/module/middleware"
	repository "github.com/labovector/vecsys-api/module/repository/event"
)

func EventRoute(adminRoute, globalRoute fiber.Router) {
	eventRepo := repository.NewEventRepositositoryImpl()
	eventController := controller.NewEventController(eventRepo)

	//	Create Event
	adminRoute.Post("/", middleware.AdminMiddleware(), eventController.CreateEvent)

	//	Edit Event
	adminRoute.Patch("/:id", middleware.AdminMiddleware(), eventController.UpdateEvent)

	//	Get All Event
	adminRoute.Get("/", middleware.AdminMiddleware(), eventController.GetAllEvent)

	// Get Event By ID
	globalRoute.Get("/event/:id", middleware.UserMiddleware(), eventController.GetEventById)

	//	Toggle Event Status
	adminRoute.Put("/:id", middleware.AdminMiddleware(), eventController.ToggleEventActive)

	//	Delete Event
	adminRoute.Delete("/:id", middleware.AdminMiddleware(), eventController.DeleteEvent)
}
