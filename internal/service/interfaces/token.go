package service

import (
	"context"

	"github.com/rafin007/api-gateway/internal/api/handler/response"
	"github.com/rafin007/api-gateway/internal/models"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *models.User) (*response.AccessToken, error)
	VerifyAccessToken(ctx context.Context, token string) (*response.TokenClaims, error)
	GenerateRefreshToken(ctx context.Context, user *models.User) (*models.RefreshToken, error)
}
