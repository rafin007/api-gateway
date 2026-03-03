package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rafin007/api-gateway/errors"
	"github.com/rafin007/api-gateway/internal/config"
	repoMocks "github.com/rafin007/api-gateway/internal/mocks/repository"
	"github.com/rafin007/api-gateway/internal/models"
	"github.com/rafin007/api-gateway/internal/service"
	serviceInterface "github.com/rafin007/api-gateway/internal/service/interfaces"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTokenService(t *testing.T, expiryTime string, signingSecret string) serviceInterface.TokenService {
	t.Helper()
	mockTokenRepo := new(repoMocks.TokenRepository)
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		AccessTokenExpiryTime: expiryTime,
		SigningSecret:         signingSecret,
	}
	return service.NewTokenService(mockTokenRepo, logger.Sugar(), cfg)
}

var testUser = &models.User{
	ID:    1,
	Email: "test@example.com",
}

// --- GenerateAccessToken ---

func TestGenerateAccessToken_Success(t *testing.T) {
	svc := setupTokenService(t, "15", "supersecret")

	result, err := svc.GenerateAccessToken(context.Background(), testUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// token string must be non-empty
	assert.NotEmpty(t, result.AccessToken)
	// token ID (jti) must be a non-empty UUID
	assert.NotEmpty(t, result.AccessTokenID)
}

func TestGenerateAccessToken_TokenContainsCorrectClaims(t *testing.T) {
	svc := setupTokenService(t, "15", "supersecret")

	result, err := svc.GenerateAccessToken(context.Background(), testUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// parse the token without verifying to inspect claims
	token, _, parseErr := jwt.NewParser().ParseUnverified(result.AccessToken, jwt.MapClaims{})
	assert.NoError(t, parseErr)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	// claims must carry the correct user identity
	assert.Equal(t, float64(testUser.ID), claims["UserID"])
	assert.Equal(t, testUser.Email, claims["Email"])

	// jti must match the returned AccessTokenID
	assert.Equal(t, result.AccessTokenID, claims["jti"])

	// token must expire in the future
	exp := int64(claims["exp"].(float64))
	assert.Greater(t, exp, time.Now().Unix())
}

func TestGenerateAccessToken_InvalidExpiryConfig(t *testing.T) {
	// AccessTokenExpiryTime is not a number — strconv.Atoi will fail
	svc := setupTokenService(t, "not-a-number", "supersecret")

	result, err := svc.GenerateAccessToken(context.Background(), testUser)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrInternalServerError)
}

func TestGenerateAccessToken_DifferentUsersProduceDifferentTokens(t *testing.T) {
	svc := setupTokenService(t, "15", "supersecret")

	user2 := &models.User{ID: 2, Email: "other@example.com"}

	result1, err1 := svc.GenerateAccessToken(context.Background(), testUser)
	result2, err2 := svc.GenerateAccessToken(context.Background(), user2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// two different users must not get the same token
	assert.NotEqual(t, result1.AccessToken, result2.AccessToken)
	assert.NotEqual(t, result1.AccessTokenID, result2.AccessTokenID)
}

// --- VerifyAccessToken ---

func TestVerifyAccessToken_Success(t *testing.T) {
	svc := setupTokenService(t, "15", "supersecret")

	generated, err := svc.GenerateAccessToken(context.Background(), testUser)
	assert.NoError(t, err)

	claims, err := svc.VerifyAccessToken(context.Background(), generated.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	// claims must reflect the original user
	assert.Equal(t, testUser.ID, claims.UserID)
	assert.Equal(t, testUser.Email, claims.Email)
}

func TestVerifyAccessToken_ExpiredToken(t *testing.T) {
	svc := setupTokenService(t, "15", "supersecret")

	// manually craft an already-expired token
	type testClaims struct {
		UserID int64
		Email  string
		jwt.RegisteredClaims
	}
	expiredClaims := &testClaims{
		UserID: testUser.ID,
		Email:  testUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenStr, _ := token.SignedString([]byte("supersecret"))

	result, err := svc.VerifyAccessToken(context.Background(), tokenStr)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
}

func TestVerifyAccessToken_WrongSigningKey(t *testing.T) {
	// generate with one secret, verify with another
	svcGenerator := setupTokenService(t, "15", "original-secret")
	svcVerifier := setupTokenService(t, "15", "different-secret")

	generated, err := svcGenerator.GenerateAccessToken(context.Background(), testUser)
	assert.NoError(t, err)

	result, err := svcVerifier.VerifyAccessToken(context.Background(), generated.AccessToken)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
}

func TestVerifyAccessToken_MalformedToken(t *testing.T) {
	svc := setupTokenService(t, "15", "supersecret")

	result, err := svc.VerifyAccessToken(context.Background(), "this.is.not.a.jwt")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
}

func TestVerifyAccessToken_EmptyToken(t *testing.T) {
	svc := setupTokenService(t, "15", "supersecret")

	result, err := svc.VerifyAccessToken(context.Background(), "")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errors.ErrUnauthorized)
}
