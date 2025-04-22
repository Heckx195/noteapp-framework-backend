package config

import (
	"fmt"
	"log"

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

	fmt.Println("Connected to PostgreSQL and migrated models!")
}

// Info:
// migrate -path db/migrations -database "postgres://user:password@localhost:5432/mydatabase?sslmode=disable" up
