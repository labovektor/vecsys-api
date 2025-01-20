package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/controller"
	"github.com/labovector/vecsys-api/module/middleware"
	repository "github.com/labovector/vecsys-api/module/repository/admin"
)

func AdminRoute(adminRoute fiber.Router) {
	adminRepo := repository.NewAdminRepositoryImpl()
	adminController := controller.NewAdminController(adminRepo)

	// Admin Auth Route
	adminRoute.Get("/", middleware.AdminMiddleware(), adminController.GetAdmin)
	adminRoute.Patch("/", middleware.AdminMiddleware(), adminController.UpdateAdminProfile)

}
