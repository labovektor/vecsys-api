package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/labovector/vecsys-api/module/controller"
	"github.com/labovector/vecsys-api/module/middleware"
	"github.com/labovector/vecsys-api/module/repository"
)

func UserRoute(userRoute fiber.Router, ctx context.Context) {
	userRepo := repository.NewUserRepositoryImpl()
	userController := controller.NewUserController(userRepo)

	// Admin Auth Route
	userRoute.Get("/", middleware.UserMiddleware(), userController.GetUser)
	// userRoute.Put("/", middleware.AdminMiddleware(), adminController.UpdateAdminProfile)

}
