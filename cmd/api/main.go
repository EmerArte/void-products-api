package main

import (
	"log"

	"github.com/emerarteaga/products-api/internal/app"
	"github.com/emerarteaga/products-api/internal/config"
	"github.com/emerarteaga/products-api/internal/infra/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present (development)
	// In production, environment variables should be set by the platform
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.Logger.Level, cfg.Logger.Format)
	logger.Info("application starting", "version", "1.0.0", "mode", cfg.Server.Mode)

	// Create and start server
	server := app.NewServer(cfg)
	if err := server.Start(); err != nil {
		logger.Error("failed to start server", "error", err)
		log.Fatal(err)
	}
}
