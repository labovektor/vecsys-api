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
	database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		panic("Cannot connect to database")
	}

	// TODO: Set Migration
	database.AutoMigrate(
		&entity.Admin{},
		&entity.Event{},
		&entity.Voucher{},
		&entity.PaymentOption{},
		&entity.Region{},
		&entity.Category{},
		&entity.Institution{},
		&entity.Payment{},
		&entity.Participant{},
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
