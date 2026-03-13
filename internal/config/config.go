package config

import "os"

type Config struct {
	Port        string
	AuthToken   string
	ReadTimeout int
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		AuthToken:   getEnv("AUTH_TOKEN", "default_token_change_me"),
		ReadTimeout: 10,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
