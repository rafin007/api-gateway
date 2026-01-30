package repository

import (
	"context"

	"github.com/rafin007/api-gateway/internal/models"
)

type TokenRepository interface {
	// GenerateAccessToken(user *models.User) (response.AccessToken, error)
	// VerifyAccessToken(token string) (bool, error)
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
}
