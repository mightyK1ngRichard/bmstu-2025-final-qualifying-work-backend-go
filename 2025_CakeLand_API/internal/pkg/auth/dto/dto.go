package dto

import (
	"github.com/google/uuid"
	"time"
)

// Register

type RegisterReq struct {
	Email       string
	Password    string
	Fingerprint string
}

type RegisterRes struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Time
}

// Login

type LoginReq struct {
	Email       string
	Password    string
	Fingerprint string
}

type LoginRes struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Time
}

// UpdateAccessToken

type UpdateAccessTokenReq struct {
	RefreshToken string
	Fingerprint  string
}

type UpdateAccessTokenRes struct {
	AccessToken string
	ExpiresIn   time.Time
}

// Logout

type LogoutReq struct {
	RefreshToken string
	Fingerprint  string
}

type LogoutRes struct {
	Message string
}

// CreateUser

type CreateUserReq struct {
	UUID             uuid.UUID
	Email            string
	PasswordHash     []byte
	RefreshTokensMap map[string]string
}

// GetUserByEmail

type GetUserByEmailReq struct {
	Email string
}

type GetUserByEmailRes struct {
	ID               uuid.UUID
	Email            string
	PasswordHash     []byte
	RefreshTokensMap map[string]string // key: fingerprint, value: refreshToken
}

// UpdateUserRefreshTokens

type UpdateUserRefreshTokensReq struct {
	UserID           string
	RefreshTokensMap map[string]string
}

// GetUserRefreshTokens

type GetUserRefreshTokensReq struct {
	UserID string
}

type GetUserRefreshTokensRes struct {
	RefreshTokensMap map[string]string
}
