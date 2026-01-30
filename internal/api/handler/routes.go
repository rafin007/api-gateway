package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rafin007/api-gateway/internal/api/middleware"
	service "github.com/rafin007/api-gateway/internal/service/interfaces"
	"go.uber.org/zap"
)

type RouterConfig struct {
	Router       *gin.Engine
	UserService  service.UserService
	TokenService service.TokenService
	Logger       *zap.SugaredLogger
}

func SetupRoutes(rc *RouterConfig) {
	r := rc.Router
	r.Use(middleware.ErrorHandler(rc.Logger))

	v1 := r.Group("/api/v1")
	{
		userHandler := NewUserHandler(rc.UserService, rc.Logger)
		v1.POST("/users/register", userHandler.RegisterUser)
		v1.POST("/users/login", userHandler.LoginUser)
	}
}
