package config

import (
	"log"
	"os"
)

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}
	return secret
}
