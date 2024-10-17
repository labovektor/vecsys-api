package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/controller"
	"github.com/labovector/vecsys-api/module/middleware"
	"github.com/labovector/vecsys-api/module/repository"
)

func AuthRoute(adminRoute fiber.Router, userRoute fiber.Router, ctx context.Context) {
	adminRepo := repository.NewAdminRepositoryImpl()
	userRepo := repository.NewUserRepositoryImpl()
	authController := controller.NewAuthController(adminRepo, userRepo)

	// Admin Auth Route
	adminRoute.Get("/", middleware.AdminMiddleware(), authController.GetAdmin)
	adminRoute.Post("/signup", authController.RegisterAdmin)
	adminRoute.Post("/login", authController.LoginAdmin)
	adminRoute.Get("/logout", middleware.AdminMiddleware(), authController.LogoutAdmin)

	userRoute.Get("/", middleware.UserMiddleware(), authController.GetUser)
	userRoute.Post("/signup", authController.RegisterUser)
	userRoute.Post("/login", authController.LoginUser)
	userRoute.Get("/logout", middleware.UserMiddleware(), authController.LogoutUser)
}
