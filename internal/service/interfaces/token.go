package service

import (
	"context"

	"github.com/rafin007/api-gateway/internal/api/handler/response"
	"github.com/rafin007/api-gateway/internal/models"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *models.User) (*response.AccessToken, error)
	// VerifyAccessToken(token string) (bool, error)
	GenerateRefreshToken(ctx context.Context, user *models.User) (*models.RefreshToken, error)
}
