package service

import (
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rafin007/api-gateway/errors"
	"github.com/rafin007/api-gateway/internal/api/handler/response"
	"github.com/rafin007/api-gateway/internal/config"
	"github.com/rafin007/api-gateway/internal/models"
	repository "github.com/rafin007/api-gateway/internal/repository/interfaces"
	service "github.com/rafin007/api-gateway/internal/service/interfaces"
	"go.uber.org/zap"
)

type tokenService struct {
	tokenRepo repository.TokenRepository
	logger    *zap.SugaredLogger
	config    *config.Config
}

func NewTokenService(tokenRepo repository.TokenRepository, logger *zap.SugaredLogger, config *config.Config) service.TokenService {
	return &tokenService{
		tokenRepo: tokenRepo,
		logger:    logger,
		config:    config,
	}
}

type claims struct {
	UserID int64
	Email  string
	jwt.RegisteredClaims
}

func (s *tokenService) GenerateAccessToken(ctx context.Context, user *models.User) (*response.AccessToken, error) {
	tokenID := uuid.NewString()
	expiryMinutes, err := strconv.Atoi(s.config.AccessTokenExpiryTime)
	if err != nil {
		s.logger.Errorw("Error converting Access token expiry time to integer")
		return nil, errors.ErrInternalServerError
	}
	expiryTime := jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMinutes) * time.Minute))

	claims := &claims{
		user.ID,
		user.Email,
		jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: expiryTime,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.config.SigningSecret))
	if err != nil {
		s.logger.Errorw("Unable to sign the token", "error", err.Error())
		return &response.AccessToken{}, errors.ErrInternalServerError
	}

	response := &response.AccessToken{
		AccessTokenID: tokenID,
		AccessToken:   tokenStr,
	}
	return response, nil
}

func (s *tokenService) GenerateRefreshToken(ctx context.Context, user *models.User) (*models.RefreshToken, error) {
	return nil, nil
}
