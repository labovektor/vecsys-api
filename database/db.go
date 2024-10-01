package database

import (
	"os"
	"time"

	"github.com/labovector/vecsys-api/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connnect to postgresql database
func InitDB() {
	conn_url := os.Getenv("DB_URL")

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
	DB = database
}
