package route

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoute(r *fiber.App, ctx context.Context) {
	adminRoutes := r.Group("/admin")
	AuthRoute(adminRoutes, ctx)
}
