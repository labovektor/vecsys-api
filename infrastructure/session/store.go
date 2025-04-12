package session

import (
	"runtime"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/labovector/vecsys-api/infrastructure/config"
)

func InitStore(config *config.RedisConfig) *session.Store {
	storage := redis.New(redis.Config{
		Host:      config.Host,
		Port:      config.Port,
		Username:  config.Username,
		Password:  config.Password,
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	store := session.New(session.Config{
		Storage:        storage,
		CookieSameSite: "None",
		CookieHTTPOnly: true,
	})

	return store
}
