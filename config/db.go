package config

import (
	"fmt"
	"log"
	"noteapp-framework-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInit() {
	var err error

	dsn := "host=localhost user=user password=password dbname=mydatabase port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Auto-migrate models -> creates the tables in the database
	err = DB.AutoMigrate(&models.Note{}, &models.Notebook{})
	if err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	fmt.Println("Connected to PostgreSQL and migrated models!")
}
