package service_test

import (
	"context"
	"testing"

	"github.com/rafin007/api-gateway/errors"
	"github.com/rafin007/api-gateway/internal/api/handler/request"
	"github.com/rafin007/api-gateway/internal/api/handler/response"
	repoMocks "github.com/rafin007/api-gateway/internal/mocks/repository"
	serviceMocks "github.com/rafin007/api-gateway/internal/mocks/service"
	"github.com/rafin007/api-gateway/internal/models"
	"github.com/rafin007/api-gateway/internal/service"
	serviceInterface "github.com/rafin007/api-gateway/internal/service/interfaces"
	"github.com/rafin007/api-gateway/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupUserService(t *testing.T) (serviceInterface.UserService, *repoMocks.UserRepository, *serviceMocks.TokenService) {
	t.Helper()
	mockUserRepo := new(repoMocks.UserRepository)
	mockTokenService := new(serviceMocks.TokenService)
	logger, _ := zap.NewDevelopment()
	svc := service.NewUserService(mockUserRepo, logger.Sugar(), mockTokenService)
	return svc, mockUserRepo, mockTokenService
}

// --- RegisterUser ---

func TestRegisterUser_Success(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// user is mutated (PasswordHash set) before CreateUser is called, so use mock.Anything
	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "test@example.com" && u.PasswordHash != ""
	})).Return(nil)

	mockTokenService.On("GenerateAccessToken", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "test@example.com"
	})).Return(&response.AccessToken{
		AccessToken: "mocked.jwt.token",
	}, nil)

	result, err := svc.RegisterUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mocked.jwt.token", result.AccessToken)
	// password hash must be populated
	assert.NotEmpty(t, user.PasswordHash)
	// plain password should NOT be stored as hash
	assert.NotEqual(t, "password1234", user.PasswordHash)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	user := &models.User{
		Email:    "existing@example.com",
		Password: "password1234",
	}

	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "existing@example.com"
	})).Return(errors.ErrUserAlreadyExists)

	result, err := svc.RegisterUser(context.Background(), user)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUserAlreadyExists)
	// token service must NOT be called when user already exists
	mockTokenService.AssertNotCalled(t, "GenerateAccessToken")
	mockUserRepo.AssertExpectations(t)
}

func TestRegisterUser_CreateUserInternalError(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password1234",
	}

	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.ErrInternalServerError)

	result, err := svc.RegisterUser(context.Background(), user)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrInternalServerError)
	mockTokenService.AssertNotCalled(t, "GenerateAccessToken")
	mockUserRepo.AssertExpectations(t)
}

func TestRegisterUser_TokenGenerationFails(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password1234",
	}

	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	mockTokenService.On("GenerateAccessToken", mock.Anything, mock.Anything).Return(nil, errors.ErrInternalServerError)

	result, err := svc.RegisterUser(context.Background(), user)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrInternalServerError)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

// --- LoginUser ---

func TestLoginUser_Success(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	plainPassword := "password1234"
	hashedPassword, _ := utils_generateHash(plainPassword) // use your utils.GenerateHashFromPassword

	storedUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	reqUser := &request.UserLogin{
		Email:    "test@example.com",
		Password: plainPassword,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)
	mockTokenService.On("GenerateAccessToken", mock.Anything, storedUser).Return(&response.AccessToken{
		AccessToken: "mocked.jwt.token",
	}, nil)

	result, err := svc.LoginUser(context.Background(), reqUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mocked.jwt.token", result.AccessToken.AccessToken)
	assert.Equal(t, "test@example.com", result.User.Email)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	reqUser := &request.UserLogin{
		Email:    "notfound@example.com",
		Password: "password1234",
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(&models.User{}, errors.ErrUserNotFound)

	result, err := svc.LoginUser(context.Background(), reqUser)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)
	mockTokenService.AssertNotCalled(t, "GenerateAccessToken")
	mockUserRepo.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	plainPassword := "correctpassword"
	hashedPassword, _ := utils_generateHash(plainPassword)

	storedUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	reqUser := &request.UserLogin{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)

	result, err := svc.LoginUser(context.Background(), reqUser)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
	mockTokenService.AssertNotCalled(t, "GenerateAccessToken")
	mockUserRepo.AssertExpectations(t)
}

func TestLoginUser_TokenGenerationFails(t *testing.T) {
	svc, mockUserRepo, mockTokenService := setupUserService(t)

	plainPassword := "password1234"
	hashedPassword, _ := utils_generateHash(plainPassword)

	storedUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	reqUser := &request.UserLogin{
		Email:    "test@example.com",
		Password: plainPassword,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)
	mockTokenService.On("GenerateAccessToken", mock.Anything, storedUser).Return(nil, errors.ErrInternalServerError)

	result, err := svc.LoginUser(context.Background(), reqUser)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrInternalServerError)
	mockUserRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

// --- GetUserByEmail ---

func TestGetUserByEmail_Success(t *testing.T) {
	svc, mockUserRepo, _ := setupUserService(t)

	expected := &models.User{ID: 1, Email: "test@example.com"}
	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expected, nil)

	result, err := svc.GetUserByEmail(context.Background(), "test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	svc, mockUserRepo, _ := setupUserService(t)

	mockUserRepo.On("GetByEmail", mock.Anything, "ghost@example.com").Return(&models.User{}, errors.ErrUserNotFound)

	result, err := svc.GetUserByEmail(context.Background(), "ghost@example.com")

	assert.ErrorIs(t, err, errors.ErrUserNotFound)
	assert.Equal(t, &models.User{}, result)
	mockUserRepo.AssertExpectations(t)
}

// --- GetUserByID ---

func TestGetUserByID_Success(t *testing.T) {
	svc, mockUserRepo, _ := setupUserService(t)

	expected := &models.User{ID: 42, Email: "test@example.com"}
	mockUserRepo.On("GetByID", mock.Anything, int64(42)).Return(expected, nil)

	result, err := svc.GetUserByID(context.Background(), 42)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	svc, mockUserRepo, _ := setupUserService(t)

	mockUserRepo.On("GetByID", mock.Anything, int64(999)).Return(&models.User{}, errors.ErrUserNotFound)

	result, err := svc.GetUserByID(context.Background(), 999)

	assert.ErrorIs(t, err, errors.ErrUserNotFound)
	assert.Equal(t, &models.User{}, result)
	mockUserRepo.AssertExpectations(t)
}

// helper to call the real hash function for test setup
// replace with: utils.GenerateHashFromPassword if exported
func utils_generateHash(password string) (string, error) {
	// import and call your actual utils.GenerateHashFromPassword here
	// e.g.: return utils.GenerateHashFromPassword(password)
	return utils.GenerateHashFromPassword(password)
}
