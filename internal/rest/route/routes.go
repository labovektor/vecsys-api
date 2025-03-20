package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/internal/rest/controller"
	"github.com/labovector/vecsys-api/internal/rest/middleware"
	ar "github.com/labovector/vecsys-api/internal/rest/repository/admin"
	er "github.com/labovector/vecsys-api/internal/rest/repository/event"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
)

type AllRepository struct {
	AdminRepository ar.AdminRepository
	UserRepository  ur.UserRepository
	EventRepository er.EventRepository
}

type AllController struct {
	AdminController *controller.AdminController
	UserController  *controller.UserController
	EventController *controller.EventController
	AuthController  *controller.AuthController
}

func SetupRoute(app *fiber.App, allRepository *AllRepository) {

	// Base route for API versioning (v1)
	api := app.Group("/api/v1")

	adminRoutes := api.Group("/admin")
	userRoutes := api.Group("/user")
	globalRoutes := api.Group("/")

	// Initialize all controllers
	allController := AllController{
		AdminController: controller.NewAdminController(allRepository.AdminRepository),
		UserController:  controller.NewUserController(allRepository.UserRepository),
		EventController: controller.NewEventController(allRepository.EventRepository),
		AuthController:  controller.NewAuthController(allRepository.AdminRepository, allRepository.UserRepository),
	}

	// Admin Auth Route
	adminAuth := adminRoutes.Group("/")
	adminAuth.Post("/signup", allController.AuthController.RegisterAdmin)
	adminAuth.Post("/login", allController.AuthController.LoginAdmin)
	adminAuth.Get("/logout", middleware.AdminMiddleware(), allController.AuthController.LogoutAdmin)

	// User Auth Route
	userAuth := userRoutes.Group("/")
	userAuth.Post("/signup", allController.AuthController.RegisterUser)
	userAuth.Post("/login", allController.AuthController.LoginUser)
	userAuth.Get("/logout", middleware.UserMiddleware(), allController.AuthController.LogoutUser)

	// Admin Route
	adminRoutes.Get("/", middleware.AdminMiddleware(), allController.AdminController.GetAdmin)
	adminRoutes.Patch("/", middleware.AdminMiddleware(), allController.AdminController.UpdateAdminProfile)

	// User Route
	userRoutes.Get("/", middleware.UserMiddleware(), allController.UserController.GetUser)
	// userRoutes.Patch("/", middleware.UserMiddleware(), allController.UserController.)

	// Event Route
	adminRoutes.Post("/", middleware.AdminMiddleware(), allController.EventController.CreateEvent)
	adminRoutes.Patch("/:id", middleware.AdminMiddleware(), allController.EventController.UpdateEvent)
	adminRoutes.Get("/", middleware.AdminMiddleware(), allController.EventController.GetAllEvent)
	globalRoutes.Get("/event/:id", middleware.UserMiddleware(), allController.EventController.GetEventById)
	adminRoutes.Put("/:id", middleware.AdminMiddleware(), allController.EventController.ToggleEventActive)
	adminRoutes.Delete("/:id", middleware.AdminMiddleware(), allController.EventController.DeleteEvent)
}
