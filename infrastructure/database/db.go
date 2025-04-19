package database

import (
	"time"

	"github.com/labovector/vecsys-api/entity"
	"github.com/labovector/vecsys-api/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connnect to postgresql database
func InitDB(config *config.PostgresConfig) *gorm.DB {

	conn_url := config.ConnUrl

	database, err := gorm.Open(postgres.Open(conn_url), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to database")
	}

	database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	database.Exec(`SET TIME ZONE 'Asia/Jakarta';`)
	database.Exec(`CREATE TYPE participant_progress AS ENUM ('registered', 'categorized', 'paid', 'validated', 'select_institution', 'fill_biodatas', 'locked');`)

	// TODO: Set Migration
	database.AutoMigrate(
		&entity.Admin{},
		&entity.Event{},
		&entity.Referal{},
		&entity.PaymentOption{},
		&entity.Region{},
		&entity.Category{},
		&entity.Institution{},
		&entity.Participant{},
		&entity.Payment{},
		&entity.Biodata{},
	)

	db, err := database.DB()
	if err != nil {
		panic("Cant connect to database")
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return database
}
