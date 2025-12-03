package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	CORS     CORSConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port int
	Mode string // debug, release, test
}

// CORSConfig holds CORS-specific configuration
type CORSConfig struct {
	AllowedOrigins []string // List of allowed origins (e.g., "http://localhost:3000")
	AllowedMethods []string // List of allowed HTTP methods
	AllowedHeaders []string // List of allowed headers
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	URI         string
	Name        string
	MaxPoolSize uint64
	Timeout     int // in seconds
}

// LoggerConfig holds logger-specific configuration
type LoggerConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, text
}

// LoadConfig reads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			URI:         getEnv("DATABASE_URI", "mongodb://localhost:27017"),
			Name:        getEnv("DATABASE_NAME", "products_db"),
			MaxPoolSize: getEnvAsUint64("DATABASE_MAX_POOL_SIZE", 100),
			Timeout:     getEnvAsInt("DATABASE_TIMEOUT", 10),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOGGER_LEVEL", "info"),
			Format: getEnv("LOGGER_FORMAT", "json"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "X-Requested-With"}),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt reads an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsUint64 reads an environment variable as uint64 or returns a default value
func getEnvAsUint64(key string, defaultValue uint64) uint64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsSlice reads an environment variable as comma-separated list or returns a default value
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	// Split by comma and trim spaces
	var result []string
	for _, v := range splitAndTrim(valueStr, ",") {
		if v != "" {
			result = append(result, v)
		}
	}
	if len(result) == 0 {
		return defaultValue
	}
	return result
}

// splitAndTrim splits a string by separator and trims spaces
func splitAndTrim(s string, sep string) []string {
	var result []string
	for _, part := range split(s, sep) {
		trimmed := trim(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// split splits a string by separator
func split(s string, sep string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

// trim removes leading and trailing spaces
func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.URI == "" {
		return fmt.Errorf("database URI is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	validModes := map[string]bool{"debug": true, "release": true, "test": true}
	if !validModes[c.Server.Mode] {
		return fmt.Errorf("invalid server mode: %s", c.Server.Mode)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Logger.Level] {
		return fmt.Errorf("invalid logger level: %s", c.Logger.Level)
	}

	return nil
}
