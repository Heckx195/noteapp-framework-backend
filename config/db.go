package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

// DBInit initializes the PostgreSQL database connection
func DBInit() {
	var err error

	connStr := "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"

	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")
}
