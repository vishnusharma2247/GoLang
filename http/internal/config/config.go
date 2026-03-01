package config

import "os"

type Config struct {
	AppEnv      string
	AppAddr     string
	DatabaseURL string
}

func Load() Config {
	return Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		AppAddr:     getEnv("APP_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
