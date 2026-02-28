package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost                 string `mapstructure:"DB_HOST" validate:"required"`
	DBName                 string `mapstructure:"DB_NAME" validate:"required"`
	DBPassword             string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBPort                 string `mapstructure:"DB_PORT" validate:"required"`
	DBUser                 string `mapstructure:"DB_USER" validate:"required"`
	Port                   string `mapstructure:"PORT" validate:"required"`
	AccessTokenExpiryTime  string `mapstructure:"ACCESS_TOKEN_EXPIRY_TIME" validate:"required"`
	RefreshTokenExpiryTime string `mapstructure:"REFRESH_TOKEN_EXPIRY_TIME" validate:"required"`
	SigningSecret          string `mapstructure:"SIGNING_SECRET" validate:"required"`
}

func LoadConfig(path string, file string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(file)

	// Tell Viper to replace dots and hyphens with underscores when matching
	// env vars, e.g. DB.HOST → DB_HOST
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	// Explicitly bind each env var to its mapstructure key.
	// This is required because Viper only auto-resolves env vars for keys
	// it has already seen — BindEnv registers them unconditionally.
	envVars := []string{
		"DB_HOST",
		"DB_NAME",
		"DB_PASSWORD",
		"DB_PORT",
		"DB_USER",
		"PORT",
		"ACCESS_TOKEN_EXPIRY_TIME",
		"REFRESH_TOKEN_EXPIRY_TIME",
		"SIGNING_SECRET",
	}
	for _, key := range envVars {
		if err := viper.BindEnv(key); err != nil {
			return config, err
		}
	}

	err = viper.ReadInConfig()

	// if the error is other than file not found, then return the error
	if err != nil && !os.IsNotExist(err) {
		return config, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return config, err
	}

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
