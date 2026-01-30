package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafin007/api-gateway/internal/api/handler"
	"github.com/rafin007/api-gateway/internal/config"
	"github.com/rafin007/api-gateway/internal/repository"
	"github.com/rafin007/api-gateway/internal/service"
	"go.uber.org/zap"
)

func Start(config *config.Config, logger *zap.SugaredLogger, pool *pgxpool.Pool) error {
	// define repositories
	userRepo := repository.NewUserRepository(pool, logger)
	tokenRepo := repository.NewTokenRepository(pool, logger)

	// define services
	tokenService := service.NewTokenService(tokenRepo, logger, config)
	userService := service.NewUserService(userRepo, logger, tokenService)

	router := gin.Default()
	handler.SetupRoutes(&handler.RouterConfig{
		Router:       router,
		UserService:  userService,
		TokenService: tokenService,
		Logger:       logger,
	})

	if err := router.Run(":" + config.Port); err != nil {
		logger.Errorw("Error starting server", "error", err.Error())
		return err
	}

	logger.Infow("Server started at:", "port", config.Port)

	return nil
}
