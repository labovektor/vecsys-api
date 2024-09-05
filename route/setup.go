package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func SetupRoute(r *fiber.App, ctx context.Context) {
	adminRoutes := r.Group("/admin")
	userRoutes := r.Group("/user")

	AuthRoute(adminRoutes, userRoutes, ctx)
}
