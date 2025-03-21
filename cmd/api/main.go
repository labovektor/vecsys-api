package main

import (
	"log"
	"os"

	"github.com/labovector/vecsys-api/infrastructure/config"
	"github.com/labovector/vecsys-api/infrastructure/database"
	"github.com/labovector/vecsys-api/infrastructure/session"
	"github.com/labovector/vecsys-api/internal/rest"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize Database
	db := database.InitDB(&conf.Postgres)

	// Initialize Store
	session := session.InitStore(&conf.Redis)

	// Custom File Writer for logger
	file, err := os.OpenFile("../../vecsys-logger.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	app := rest.New(session, db, file)

	// Run the app
	if err := app.Listen(":8787"); err != nil {
		panic(err)
	}
}
