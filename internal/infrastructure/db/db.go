package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"wati-clone-backend/internal/core/models"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for local dev if env not set
		dsn = "host=localhost user=postgres password=postgres dbname=wati_clone port=5432 sslmode=disable"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Auto Migrate
	err = DB.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Contact{},
		&models.Message{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
