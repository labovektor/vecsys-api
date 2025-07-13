package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/infrastructure/email"
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
	ReferalController                   *controller.ReferalController
	ParticipantAdministrationController *controller.ParticipantAdministrationController
	ParticipantDataController           *controller.ParticipantDataController
	PaymentOptionController             controller.PaymentOptionController
	InstituionController                controller.InstitutionController
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
		RegionController:  controller.NewRegionController(allRepository.RegionRepository),
		ReferalController: controller.NewReferalController(allRepository.ReferalRepository),
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
		PaymentOptionController: controller.NewPaymentOptionController(
			allRepository.PaymentRepository,
		),
		InstituionController: *controller.NewInstitutionController(
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
	adminRoutes.Get("/event/:id/category", middleware.AdminMiddleware(), allController.CategoryController.GetAllCategoryByEventId)
	adminRoutes.Get("/category/:id", middleware.AdminMiddleware(), allController.CategoryController.GetCategoryById)
	adminRoutes.Post("/event/:id/category", middleware.AdminMiddleware(), allController.CategoryController.AddCategoryToEvent)
	adminRoutes.Patch("/category/:id", middleware.AdminMiddleware(), allController.CategoryController.UpdateCategory)
	adminRoutes.Delete("/category/:id", middleware.AdminMiddleware(), allController.CategoryController.DeleteCategory)

	// Event Region Route
	event.Get("/:id/region", allController.RegionController.GetAllRegionsByEventId)
	adminRoutes.Get("/region/:id", middleware.AdminMiddleware(), allController.RegionController.GetRegionById)
	event.Post("/:id/region", allController.RegionController.AddRegionToEvent)
	adminRoutes.Patch("/region/:id", middleware.AdminMiddleware(), allController.RegionController.UpdateRegion)
	adminRoutes.Delete("/region/:id", middleware.AdminMiddleware(), allController.RegionController.DeleteRegion)

	// Event Payment Option Route
	event.Get("/:id/payment-option", allController.PaymentOptionController.GetAllPaymentOptionByEventId)
	adminRoutes.Get("/payment-option/:id", middleware.AdminMiddleware(), allController.PaymentOptionController.GetPaymentOptionById)
	event.Post("/:id/payment-option", allController.PaymentOptionController.AddPaymentOptionToEvent)
	adminRoutes.Patch("/payment-option/:id", middleware.AdminMiddleware(), allController.PaymentOptionController.UpdatePaymentOption)
	adminRoutes.Delete("/payment-option/:id", middleware.AdminMiddleware(), allController.PaymentOptionController.DeletePaymentOption)

	// Event Referal Route
	event.Get("/:id/referal", allController.ReferalController.GetReferalsByEventId)
	adminRoutes.Get("/referal/:code", middleware.AdminMiddleware(), allController.ReferalController.GetReferalByCode)
	event.Post("/:id/referal", allController.ReferalController.AddReferalToEvent)
	adminRoutes.Delete("/referal/:id", middleware.AdminMiddleware(), allController.ReferalController.DeleteReferal)

	// Event Participant Route
	event.Get("/:id/participant", allController.UserController.GetAllParticipant)
	adminRoutes.Patch("/participant/:id/verify", middleware.AdminMiddleware(), allController.UserController.VerifyParticipant)
	adminRoutes.Get("/participant/:id/card", middleware.AdminMiddleware(), allController.UserController.GeneratePdfParticipant)
	event.Post("/:id/participant/bulk", allController.UserController.BulkAddParticipantFromCSV)
	adminRoutes.Get("/participant/:id", middleware.AdminMiddleware(), allController.UserController.GetParticipantByID)
	adminRoutes.Patch("/participant/:id", middleware.AdminMiddleware(), allController.UserController.UpdateParticipantData)
	adminRoutes.Patch("/participant/:id/biodatas", middleware.AdminMiddleware(), allController.UserController.UpdateParticipantBiodata)
	adminRoutes.Delete("/participant/:id", middleware.AdminMiddleware(), allController.UserController.DeleteParticipant)

	// Event Institution Route
	event.Get("/:id/institution", allController.InstituionController.GetAllInstitutions)
	adminRoutes.Get("/institution/:id", middleware.AdminMiddleware(), allController.InstituionController.GetInstitutions)
	event.Post("/:id/institution", allController.InstituionController.AddInstitutions)
	adminRoutes.Patch("/institution/:id", middleware.AdminMiddleware(), allController.InstituionController.UpdateInstitutions)
	adminRoutes.Delete("/institution/:id", middleware.AdminMiddleware(), allController.InstituionController.DeleteInstitutions)

	// Write your route up here

	// Participant Route
	userAdministration := userRoutes.Group("/administration", middleware.UserMiddleware())
	// Get All Event Category & Region
	userAdministration.Get("/category", allController.ParticipantAdministrationController.GetAllEventCategoryAndRegion)
	// Pick Category and Region
	userAdministration.Patch("/category", middleware.UserStepMiddleware(allRepository.UserRepository, entity.StepCategorizedParticipant), allController.ParticipantAdministrationController.PickCategoryAndRegion)
	// Get All Payment Option
	userAdministration.Get("/payment", allController.ParticipantAdministrationController.GetAllPaymentOption)
	// Validate Referal
	userAdministration.Post("/validate-referal", middleware.UserStepMiddleware(allRepository.UserRepository, entity.StepPaidParticipant), allController.ParticipantAdministrationController.ValidateReferal)
	// Payment
	userAdministration.Patch("/payment", middleware.UserStepMiddleware(allRepository.UserRepository, entity.StepPaidParticipant), allController.ParticipantAdministrationController.Payment)

	userData := userRoutes.Group("/data", middleware.UserMiddleware())
	// Get All Institution
	userData.Get("/institution", allController.ParticipantDataController.GetAllInstitution)
	// Add Institution
	userData.Post("/institution", allController.ParticipantDataController.AddInstitution)
	// Pick Institution
	userData.Patch("/institution", middleware.UserStepMiddleware(allRepository.UserRepository, entity.StepSelectInstitutionParticipant), allController.ParticipantDataController.PickInstitution)
	// Add Members
	userData.Post("/member", middleware.UserStepMiddleware(allRepository.UserRepository, entity.StepFillBiodatasParticipant), allController.ParticipantDataController.AddMembers)
	// Remove Members
	userData.Delete("/member/:id", middleware.UserStepMiddleware(allRepository.UserRepository, entity.StepFillBiodatasParticipant), allController.ParticipantDataController.RemoveMembers)
	// Lock Data
	userData.Patch("/lock", middleware.UserStepMiddleware(allRepository.UserRepository, entity.StepLockedParticipant), allController.ParticipantDataController.LockData)
	// Download Card
	globalRoutes.Get("/card/:id", allController.UserController.GeneratePdfParticipant)
}
