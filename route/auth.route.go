package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/controller"
	"github.com/labovector/vecsys-api/module/middleware"
	adminRepo "github.com/labovector/vecsys-api/module/repository/admin"
	userRepo "github.com/labovector/vecsys-api/module/repository/user"
)

func AuthRoute(adminRoute fiber.Router, userRoute fiber.Router, ctx context.Context) {
	adminRepo := adminRepo.NewAdminRepositoryImpl()
	userRepo := userRepo.NewUserRepositoryImpl()
	authController := controller.NewAuthController(adminRepo, userRepo)

	// Admin Auth Route
	adminRoute.Post("/signup", authController.RegisterAdmin)
	adminRoute.Post("/login", authController.LoginAdmin)
	adminRoute.Get("/logout", middleware.AdminMiddleware(), authController.LogoutAdmin)

	userRoute.Post("/signup", authController.RegisterUser)
	userRoute.Post("/login", authController.LoginUser)
	userRoute.Get("/logout", middleware.UserMiddleware(), authController.LogoutUser)
}
