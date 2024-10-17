package database

import (
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/labovector/vecsys-api/util"
)

var Store *session.Store

func InitStore() {
	storage := redis.New(redis.Config{
		Host:      os.Getenv("REDIS_HOST"),
		Port:      util.StringToInteger(os.Getenv("REDIS_PORT"), 6379),
		Username:  os.Getenv("REDIS_USERNAME"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	store := session.New(session.Config{
		Storage:        storage,
		CookieSameSite: "None",
		CookieHTTPOnly: true,
		Expiration:     2 * time.Hour,
	})

	Store = store
}
