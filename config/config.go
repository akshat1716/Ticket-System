package config

import (
	"os"
)

type Config struct {
	JWTSecret    string
	Port         string
	DatabasePath string
}

func Load() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./ticket_system.db"
	}

	return &Config{
		JWTSecret:    jwtSecret,
		Port:         port,
		DatabasePath: dbPath,
	}
}
