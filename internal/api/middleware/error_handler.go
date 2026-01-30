package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafin007/api-gateway/errors"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			if appErr, ok := err.Err.(*errors.AppError); ok {

				// check for validation errors
				if len(appErr.Errors) > 0 {
					c.JSON(appErr.Code, gin.H{
						"errors": appErr.Errors,
					})
					return
				}

				// otherwise return the message
				c.JSON(appErr.Code, gin.H{
					"error": appErr.Message,
				})
				return
			}

			// for other generic errors that are not of type AppErr, we don't wanna send back those
			logger.Errorw("Internal server error", "path", c.Request.URL.Path, "error", err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong",
			})
		}
	}
}
