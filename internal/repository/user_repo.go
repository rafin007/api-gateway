package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	customErr "github.com/rafin007/api-gateway/errors"
	"github.com/rafin007/api-gateway/internal/models"
	repository "github.com/rafin007/api-gateway/internal/repository/interfaces"
	"go.uber.org/zap"
)

type userRepo struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewUserRepository(pool *pgxpool.Pool, logger *zap.SugaredLogger) repository.UserRepository {
	return &userRepo{
		pool:   pool,
		logger: logger,
	}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users 
			(email, password_hash, first_name, last_name, phone)
			VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	r.logger.Infow("Creating a new user in the database", "email", user.Email)

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		r.logger.Errorw("Failed to create a new user in the database", "error", err.Error())
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return customErr.ErrUserAlreadyExists
			}
		}
		return customErr.ErrInternalServerError
	}

	r.logger.Infow("User created successfully in the database", "user_id", user.ID, "email", user.Email)
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `SELECT id, email, password_hash, first_name, last_name, verified, phone, created_at, updated_at 
			 FROM users WHERE email=$1`
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Verified,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Errorw("User does not exist", "email", email, "error", err.Error())
			return &user, customErr.ErrUserNotFound
		} else {
			r.logger.Errorw("Failed to fetch user", "email", email, "error", err.Error())
			return &user, customErr.ErrInternalServerError
		}
	}

	return &user, nil
}
