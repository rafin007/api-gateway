package service

import (
	"context"

	"github.com/rafin007/api-gateway/internal/api/handler/request"
	"github.com/rafin007/api-gateway/internal/api/handler/response"
	"github.com/rafin007/api-gateway/internal/models"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *models.User) (*response.AccessToken, error)
	LoginUser(ctx context.Context, userReq *request.UserLogin) (*response.UserResponse, error)
}
