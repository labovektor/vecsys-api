package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/internal/email"
	"github.com/labovector/vecsys-api/internal/rest/controller"
	"github.com/labovector/vecsys-api/internal/rest/middleware"
	ar "github.com/labovector/vecsys-api/internal/rest/repository/admin"
	cr "github.com/labovector/vecsys-api/internal/rest/repository/category"
	er "github.com/labovector/vecsys-api/internal/rest/repository/event"
	ir "github.com/labovector/vecsys-api/internal/rest/repository/institution"
	pr "github.com/labovector/vecsys-api/internal/rest/repository/payment"
	vr "github.com/labovector/vecsys-api/internal/rest/repository/referal"
	rr "github.com/labovector/vecsys-api/internal/rest/repository/region"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
	"github.com/labovector/vecsys-api/internal/util"
)

type AllRepository struct {
	AdminRepository       ar.AdminRepository
	UserRepository        ur.UserRepository
	EventRepository       er.EventRepository
	PaymentRepository     pr.PaymentRepository
	RegionRepository      rr.RegionRepository
	ReferalRepository     vr.ReferalRepository
	CategoryRepository    cr.CategoryRepository
	InstitutionRepository ir.InstitutionRepository
}

type AllController struct {
	AdminController                     *controller.AdminController
	UserController                      *controller.UserController
	EventController                     *controller.EventController
	AuthController                      *controller.AuthController
	CategoryController                  controller.CategoryController
	RegionController                    controller.RegionController
	ParticipantAdministrationController *controller.ParticipantAdministrationController
	ParticipantDataController           *controller.ParticipantDataController
}

func SetupRoute(app *fiber.App, allRepository *AllRepository, jwtMaker util.Maker, emailDialer email.EmailDialer) {

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
		AuthController:  controller.NewAuthController(allRepository.AdminRepository, allRepository.UserRepository, jwtMaker, emailDialer),
		CategoryController: controller.NewCategoryController(
			allRepository.CategoryRepository,
		),
		RegionController: controller.NewRegionController(allRepository.RegionRepository),
		ParticipantAdministrationController: controller.NewParticipantController(
			allRepository.UserRepository,
			allRepository.CategoryRepository,
			allRepository.RegionRepository,
			allRepository.PaymentRepository,
			allRepository.ReferalRepository,
		),
		ParticipantDataController: controller.NewParticipantDataController(
			allRepository.UserRepository,
			allRepository.InstitutionRepository,
		),
	}

	// Admin Auth Route
	adminAuth := adminRoutes.Group("/")
	adminAuth.Post("/register", allController.AuthController.RegisterAdmin)
	adminAuth.Post("/login", allController.AuthController.LoginAdmin)
	adminAuth.Get("/logout", middleware.AdminMiddleware(), allController.AuthController.LogoutAdmin)

	// User Auth Route
	userAuth := userRoutes.Group("/")
	userAuth.Post("/register", allController.AuthController.RegisterUser)
	userAuth.Post("/login", allController.AuthController.LoginUser)
	userAuth.Get("/logout", middleware.UserMiddleware(), allController.AuthController.LogoutUser)
	userAuth.Post("/forgot-password", allController.AuthController.ForgotPasswordUser)
	userAuth.Post("/reset-password", allController.AuthController.ResetPasswordUser)

	// Admin Route
	adminRoutes.Get("/", middleware.AdminMiddleware(), allController.AdminController.GetAdmin)
	adminRoutes.Patch("/", middleware.AdminMiddleware(), allController.AdminController.UpdateAdminProfile)

	// User Route
	userRoutes.Get("/", middleware.UserMiddleware(), allController.UserController.GetUser)
	userRoutes.Get("/data", middleware.UserMiddleware(), allController.UserController.GetAllParticipantData)
	// Get Participant State
	userRoutes.Get("/state", middleware.UserMiddleware(), allController.UserController.GetParticipantState)

	// Event Route
	event := adminRoutes.Group("/event", middleware.AdminMiddleware())
	event.Post("/", allController.EventController.CreateEvent)
	event.Patch("/:id", allController.EventController.UpdateEvent)
	event.Get("/", allController.EventController.GetAllEvent)
	event.Put("/:id", allController.EventController.ToggleEventActive)
	event.Get("/:id/toggle", allController.EventController.ToggleEventActive)
	event.Delete("/:id", allController.EventController.DeleteEvent)
	globalRoutes.Get("/event/:id", allController.EventController.GetEventById)

	// Evert Category Route
	adminRoutes.Get("/event/:id/category", allController.CategoryController.GetAllCategoryByEventId)
	adminRoutes.Get("/category/:id", allController.CategoryController.GetCategoryById)
	adminRoutes.Post("/event/:id/category", allController.CategoryController.AddCategoryToEvent)
	adminRoutes.Patch("/category/:id", allController.CategoryController.UpdateCategory)
	adminRoutes.Delete("/category/:id", allController.CategoryController.DeleteCategory)

	// Event Region Route
	event.Get("/:id/region", allController.RegionController.GetAllRegionsByEventId)
	adminRoutes.Get("/region/:id", allController.RegionController.GetRegionById)
	event.Post("/:id/region", allController.RegionController.AddRegionToEvent)
	adminRoutes.Patch("/region/:id", allController.RegionController.UpdateRegion)
	adminRoutes.Delete("/region/:id", allController.RegionController.DeleteRegion)

	// TODO: Test All API Route Bellow
	// Participant Route
	userAdministration := userRoutes.Group("/administration", middleware.UserMiddleware())
	// Get All Event Category & Region
	userAdministration.Get("/category", allController.ParticipantAdministrationController.GetAllEventCategoryAndRegion)
	// Pick Category and Region
	userAdministration.Patch("/category", allController.ParticipantAdministrationController.PickCategoryAndRegion)
	// Get All Payment Option
	userAdministration.Get("/payment", allController.ParticipantAdministrationController.GetAllPaymentOption)
	// Validate Referal
	userAdministration.Post("/validate-referal", allController.ParticipantAdministrationController.ValidateReferal)
	// Payment
	userAdministration.Post("/payment", allController.ParticipantAdministrationController.Payment)

	userData := userRoutes.Group("/data", middleware.UserMiddleware())
	// Get All Institution
	userData.Get("/institution", allController.ParticipantDataController.GetAllInstitution)
	// Add Institution
	userData.Post("/institution", allController.ParticipantDataController.AddInstitution)
	// Pick Institution
	userData.Patch("/institution", allController.ParticipantDataController.PickInstitution)
	// Add Members
	userData.Post("/members", allController.ParticipantDataController.AddMembers)
	// Remove Members
	userData.Delete("/members", allController.ParticipantDataController.RemoveMembers)
	// Lock Data
	userData.Patch("/lock", allController.ParticipantDataController.LockData)
}
