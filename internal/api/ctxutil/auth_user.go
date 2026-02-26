package ctxutil

import (
	"github.com/gin-gonic/gin"
	"github.com/rafin007/api-gateway/errors"
	"go.uber.org/zap"
)

type AuthUser struct {
	UserID int64
	Email  string
}

func GetAuthUser(c *gin.Context, logger *zap.SugaredLogger) (*AuthUser, error) {
	userID, ok := c.Get("userID")
	if !ok {
		logger.Error("No userID found in Gin context")
		return nil, errors.ErrUnauthorized
	}

	email, ok := c.Get("email")
	if !ok {
		logger.Error("No email found in Gin context")
		return nil, errors.ErrUnauthorized
	}

	return &AuthUser{
		UserID: userID.(int64),
		Email:  email.(string),
	}, nil
}
