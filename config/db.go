package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInit() {
	var dsn string
	if os.Getenv("GO_ENV") == "test" {
		// Use test database
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("TEST_DB_HOST"),
			os.Getenv("TEST_DB_PORT"),
			os.Getenv("TEST_DB_USER"),
			os.Getenv("TEST_DB_PASSWORD"),
			os.Getenv("TEST_DB_NAME"),
		)
	} else {
		// Use development database
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	fmt.Println("Connected to PostgreSQL with dsn=", dsn)
}

// Info:
// migrate -path db/migrations -database "postgres://user:password@localhost:5432/mydatabase?sslmode=disable" up
// GO_ENV=test go test ./handlers
