package response

import "github.com/rafin007/api-gateway/internal/models"

type AccessToken struct {
	AccessTokenID string
	AccessToken   string
}

type UserResponse struct {
	models.User
	AccessToken
}
