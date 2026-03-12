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
	ENV_NAME   = "APP_ENV"
)

func main() {
	env := os.Getenv(ENV_NAME)

	if env != "production" && env != "development" {
		panic("Error: " + ENV_NAME + " must be 'production' or 'development'")
	}

	sugLog := logger.InitLogger(env)
	defer func() {
		_ = sugLog.Sync()
	}()

	sugLog.Info("Loading configuration...")
	config, err := config.LoadConfig(configPath, configFile, env)
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
