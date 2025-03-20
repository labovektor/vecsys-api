package main

import (
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

	app := rest.New(session, db)

	// Run the app
	if err := app.Listen(":8787"); err != nil {
		panic(err)
	}
}
