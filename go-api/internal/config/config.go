package config

import (
	"os"
	"strings"
)

// Config centraliza la configuración desde variables de entorno.
type Config struct {
	Port       string
	NodeAPIURL string
	JWTSecret  string
}

func Load() Config {
	return Config{
		Port:       getEnv("PORT", "8080"),
		NodeAPIURL: strings.TrimRight(getEnv("NODE_API_URL", "http://localhost:3000"), "/"),
		JWTSecret:  getEnv("JWT_SECRET", "dev-secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}