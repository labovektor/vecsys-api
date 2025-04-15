package rest

import (
	"context"
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
	cr "github.com/labovector/vecsys-api/internal/rest/repository/category"
	er "github.com/labovector/vecsys-api/internal/rest/repository/event"
	ir "github.com/labovector/vecsys-api/internal/rest/repository/institution"
	pr "github.com/labovector/vecsys-api/internal/rest/repository/payment"
	rr "github.com/labovector/vecsys-api/internal/rest/repository/region"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
	vr "github.com/labovector/vecsys-api/internal/rest/repository/voucher"
)

func New(session *session.Store, db *gorm.DB, logFile *os.File) *fiber.App {
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

	// Use Logger
	app.Use(logger.New(logger.Config{
		Format:        "${pid} ${locals:requestid} ${status} - ${method} ${path} ${error}\n",
		TimeFormat:    "02-Jan-2006",
		TimeZone:      "Asia/Jakarta",
		Output:        logFile,
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

	route.SetupRoute(app, &route.AllRepository{
		AdminRepository:       ar.NewAdminRepositoryImpl(db),
		UserRepository:        ur.NewUserRepositoryImpl(db),
		EventRepository:       er.NewEventRepositositoryImpl(db),
		PaymentRepository:     pr.NewPaymentRepositoryImpl(db),
		RegionRepository:      rr.NewRegionRepositoryImpl(db),
		VoucherRepository:     vr.NewVoucherRepositoryImpl(db),
		CategoryRepository:    cr.NewCategoryRepositoryImpl(db),
		InstitutionRepository: ir.NewInstitutionRepositoryImpl(db),
	})

	return app
}
