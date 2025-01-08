package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/controller"
	"github.com/labovector/vecsys-api/module/middleware"
	"github.com/labovector/vecsys-api/module/repository"
)

func AdminRoute(adminRoute fiber.Router, ctx context.Context) {
	adminRepo := repository.NewAdminRepositoryImpl()
	adminController := controller.NewAdminController(adminRepo)

	// Admin Auth Route
	adminRoute.Get("/", middleware.AdminMiddleware(), adminController.GetAdmin)
	adminRoute.Put("/", middleware.AdminMiddleware(), adminController.UpdateAdminProfile)

}
