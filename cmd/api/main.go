package main

import (
	"os"

	"github.com/rafin007/api-gateway/internal/app"
	"github.com/rafin007/api-gateway/internal/config"
	"github.com/rafin007/api-gateway/internal/db"
	"github.com/rafin007/api-gateway/pkg/logger"
)

const (
	configPath = "./"
	configFile = ".env"
	MODE       = "APP_MODE"
)

func main() {
	mode := os.Getenv(MODE)

	if mode != "prod" && mode != "dev" {
		panic("Error: " + MODE + " must be 'prod' or 'dev'")
	}

	sugLog := logger.InitLogger(mode)
	defer func() {
		_ = sugLog.Sync()
	}()

	sugLog.Info("Loading configuration...")
	config, err := config.LoadConfig(configPath, configFile)
	if err != nil {
		panic("Error loading configuration: " + err.Error())
	}

	pool, err := db.InitDB(&config, sugLog)
	if err != nil {
		panic("Database connection failed!")
	}
	defer pool.Close()

	if err := app.Start(&config, sugLog, pool); err != nil {
		panic("App failed to start " + err.Error())
	}
}
