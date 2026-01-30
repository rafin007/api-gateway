package service

import (
	"context"

	"github.com/rafin007/api-gateway/errors"
	"github.com/rafin007/api-gateway/internal/api/handler/request"
	"github.com/rafin007/api-gateway/internal/api/handler/response"
	"github.com/rafin007/api-gateway/internal/models"
	repository "github.com/rafin007/api-gateway/internal/repository/interfaces"
	service "github.com/rafin007/api-gateway/internal/service/interfaces"
	"github.com/rafin007/api-gateway/internal/utils"
	"go.uber.org/zap"
)

type userService struct {
	userRepo     repository.UserRepository
	logger       *zap.SugaredLogger
	tokenService service.TokenService
}

func NewUserService(userRepo repository.UserRepository, logger *zap.SugaredLogger, tokentokenService service.TokenService) service.UserService {
	return &userService{
		userRepo:     userRepo,
		logger:       logger,
		tokenService: tokentokenService,
	}
}

func (s *userService) RegisterUser(ctx context.Context, user *models.User) (*response.AccessToken, error) {
	var res *response.AccessToken

	hashedPassword, err := utils.GenerateHashFromPassword(user.Password)
	if err != nil {
		s.logger.Errorw("Error hashing password", "error", err.Error())
		return res, errors.ErrInternalServerError
	}
	user.PasswordHash = hashedPassword

	token, err := s.tokenService.GenerateAccessToken(ctx, user)
	if err != nil {
		return res, err
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return res, err
	}

	return token, nil
}

func (s *userService) LoginUser(ctx context.Context, reqUser *request.UserLogin) (*response.UserResponse, error) {
	var res *response.UserResponse
	// check if user exists
	user, err := s.userRepo.GetByEmail(ctx, reqUser.Email)
	if err != nil {
		return res, err
	}

	// hash the password and compare it against the one in database
	if match := utils.VerifyHashAndPassword(user.PasswordHash, reqUser.Password); !match {
		s.logger.Errorw("Incorrect password", "email", user.Email)
		return res, errors.ErrInvalidCredentials
	}

	// generate access token and refresh token
	token, err := s.tokenService.GenerateAccessToken(ctx, user)
	if err != nil {
		return res, err
	}

	// return the user
	res = &response.UserResponse{
		User:        *user,
		AccessToken: *token,
	}

	return res, nil
}
