package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/rafin007/api-gateway/errors"
	"github.com/rafin007/api-gateway/internal/api/handler/request"
	"github.com/rafin007/api-gateway/internal/api/handler/response"
	"github.com/rafin007/api-gateway/internal/models"
	service "github.com/rafin007/api-gateway/internal/service/interfaces"
	"go.uber.org/zap"
)

type userHandler struct {
	userService service.UserService
	logger      *zap.SugaredLogger
}

func NewUserHandler(userService service.UserService, logger *zap.SugaredLogger) *userHandler {
	return &userHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var reqBody request.UserRegistration
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		h.logger.Errorw("Error in request body", "error", err.Error())

		validationErrors := errors.ValidateFields(err)
		if len(validationErrors) > 0 {
			c.Error(errors.ValidationError(validationErrors))
			return
		}

		c.Error(errors.BadRequest(errors.ErrBadRequest.Error()))
		return
	}

	var user models.User
	if err := copier.Copy(&user, &reqBody); err != nil {
		h.logger.Errorw("Error copying request body to user", "error", err.Error())
		c.Error(errors.InternalServerError(errors.ErrInternalServerError.Error()))
		return
	}

	token, err := h.userService.RegisterUser(c, &user)
	if err != nil {
		appErr := errors.MapServiceError(err)
		c.Error(appErr)
		return
	}

	res := &response.UserResponse{
		User:        user,
		AccessToken: *token,
	}
	c.JSON(http.StatusCreated, res)
}

func (h *userHandler) LoginUser(c *gin.Context) {
	var reqBody *request.UserLogin
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		h.logger.Errorw("Error in request body", "error", err.Error())

		validationErrors := errors.ValidateFields(err)
		if len(validationErrors) > 0 {
			c.Error(errors.ValidationError(validationErrors))
			return
		}

		c.Error(errors.BadRequest(errors.ErrBadRequest.Error()))
		return
	}

	res, err := h.userService.LoginUser(c, reqBody)
	if err != nil {
		appErr := errors.MapServiceError(err)
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, res)
}
