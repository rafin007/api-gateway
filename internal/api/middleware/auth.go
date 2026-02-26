package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rafin007/api-gateway/errors"
	service "github.com/rafin007/api-gateway/internal/service/interfaces"
	"go.uber.org/zap"
)

func Auth(logger *zap.SugaredLogger, tokenService service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the auth header
		authHeader := c.Request.Header.Get("authorization")

		// parse the bearer token from the header
		var token string
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
				token = parts[1]
			}
		}

		// handle error
		if token == "" {
			logger.Error("JWT not found")
			appErr := errors.MapServiceError(errors.ErrUnauthorized)
			c.Error(appErr)
			c.Abort()
			return
		}

		// decode the jwt and fetch the user by calling token service
		claims, err := tokenService.VerifyAccessToken(c.Request.Context(), token)
		if err != nil {
			appErr := errors.MapServiceError(errors.ErrUnauthorized)
			c.Error(appErr)
			c.Abort()
			return
		}

		// set the user in the gin context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
