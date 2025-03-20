package rest

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labovector/vecsys-api/internal/rest/route"
	"gorm.io/gorm"

	ar "github.com/labovector/vecsys-api/internal/rest/repository/admin"
	er "github.com/labovector/vecsys-api/internal/rest/repository/event"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
)

func New(session *session.Store, db *gorm.DB) *fiber.App {
	_ = context.Background()
	app := fiber.New(fiber.Config{
		AppName: "vecsys",
	})

	app.Static("/public", "../../__public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 30 * time.Second,
	})

	// Custom File Writer for logger
	file, err := os.OpenFile("../../vecsys-logger.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	// Use Logger
	app.Use(logger.New(logger.Config{
		Format:        "${pid} ${locals:requestid} ${status} - ${method} ${path} ${error}\n",
		TimeFormat:    "02-Jan-2006",
		TimeZone:      "Asia/Jakarta",
		Output:        file,
		DisableColors: true,
	}))

	app.Use(func(c *fiber.Ctx) error {
		sess, err := session.Get(c)
		if err != nil {
			return err
		}
		c.Locals("session", sess)
		return c.Next()
	})

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     "http://localhost:3000, https://vecsys.vercel.app",
			AllowCredentials: true,
		},
	))

	// app.Use(encryptcookie.New(encryptcookie.Config{
	// 	Key: "YuUkkdJqEi6u8uGMU7Hn2YvF5fSranbO",
	// }))

	// Prevent client to send too many request
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        10,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("Too many request")
		},
	}))

	// TODO: Setup Route
	route.SetupRoute(app, &route.AllRepository{

		AdminRepository: ar.NewAdminRepositoryImpl(db),
		UserRepository:  ur.NewUserRepositoryImpl(db),
		EventRepository: er.NewEventRepositositoryImpl(db),
	})

	return app
}
