package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafin007/api-gateway/internal/config"
	"go.uber.org/zap"
)

func InitDB(config *config.Config, logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Errorw("Unable to parse connection string", "error", err.Error())
		return nil, err
	}

	pgxCfg.MaxConns = 10
	pgxCfg.MinConns = 2
	pgxCfg.MaxConnLifetime = time.Hour
	pgxCfg.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		logger.Errorw("Unable to create connection pool", "error", err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = pool.Ping(ctx); err != nil {
		logger.Errorw("Unable to connect to the database", "error", err.Error())
		return nil, err
	}

	logger.Info("Successfully connected to the database")
	return pool, nil
}
