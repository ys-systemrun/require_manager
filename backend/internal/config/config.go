package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ServerPort       string
	CORSAllowOrigins []string
}

// Load reads configuration from the environment, applying sensible defaults.
func Load() *Config {
	c := &Config{
		DBHost:     env("DB_HOST", "localhost"),
		DBPort:     env("DB_PORT", "5432"),
		DBUser:     env("DB_USER", "reqmgr"),
		DBPassword: env("DB_PASSWORD", "reqmgr_pass"),
		DBName:     env("DB_NAME", "reqmgr"),
		DBSSLMode:  env("DB_SSLMODE", "disable"),
		ServerPort: env("SERVER_PORT", "8080"),
	}

	origins := env("CORS_ALLOW_ORIGINS", "http://localhost:5173")
	for _, o := range strings.Split(origins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			c.CORSAllowOrigins = append(c.CORSAllowOrigins, trimmed)
		}
	}
	return c
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tokyo",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
