package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/admin/controller"
	"github.com/labovector/vecsys-api/module/admin/middleware"
	"github.com/labovector/vecsys-api/module/admin/repository"
)

func AuthRoute(r fiber.Router, ctx context.Context) {
	adminRepo := repository.NewAdminRepositoryImpl()
	adminController := controller.NewAuthController(adminRepo)
	r.Post("/signup", adminController.Register)
	r.Post("/login", adminController.Login)
	r.Get("/logout", adminController.Logout, middleware.AdminMiddleware())
}
