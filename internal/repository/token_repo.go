package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	customErr "github.com/rafin007/api-gateway/errors"
	"github.com/rafin007/api-gateway/internal/models"
	repository "github.com/rafin007/api-gateway/internal/repository/interfaces"
	"go.uber.org/zap"
)

type tokenRepo struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewTokenRepository(pool *pgxpool.Pool, logger *zap.SugaredLogger) repository.TokenRepository {
	return &tokenRepo{
		pool:   pool,
		logger: logger,
	}
}

func (r *tokenRepo) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens
			 (user_id, token_hash, device_info, ip_address, expires_at)
			 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	r.logger.Infow("Creating a new refresh token in the database", "user_id", refreshToken.UserID)

	err := r.pool.QueryRow(
		ctx,
		query,
		refreshToken.UserID,
		refreshToken.TokenHash,
		refreshToken.DeviceInfo,
		refreshToken.IPAddress,
		refreshToken.ExpiresAt,
	).Scan(&refreshToken.ID, &refreshToken.CreatedAt)

	if err != nil {
		r.logger.Errorw("Failed to create a new refresh token in the database", "user_id", refreshToken.UserID, "error", err.Error())

		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23503":
				return customErr.ErrUserNotFound
			case "23502":
				return customErr.ErrBadRequest
			}
		}
		return customErr.ErrInternalServerError
	}

	r.logger.Infow("Refresh token created successfully in the database", "refresh_token_id", refreshToken.ID, "user_id", refreshToken.UserID)

	return nil
}
