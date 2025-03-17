package auth

import (
	"2025_CakeLand_API/internal/pkg/auth/entities"
	"context"
)

// mockgen -source=internal/pkg/auth/interfaces.go -destination=internal/pkg/auth/mocks/mock_auth.go -package=mocks

type IAuthUsecase interface {
	Register(context.Context, entities.RegisterReq) (*entities.RegisterRes, error)
	Login(context.Context, entities.LoginReq) (*entities.LoginRes, error)
	Logout(context.Context, entities.LogoutReq) (*entities.LogoutRes, error)
	UpdateAccessToken(context.Context, entities.UpdateAccessTokenReq) (*entities.UpdateAccessTokenRes, error)
}

type IAuthRepository interface {
	CreateUser(context.Context, entities.CreateUserReq) error
	GetUserByEmail(context.Context, entities.GetUserByEmailReq) (*entities.GetUserByEmailRes, error)
	UpdateUserRefreshTokens(context.Context, entities.UpdateUserRefreshTokensReq) error
	GetUserRefreshTokens(context.Context, entities.GetUserRefreshTokensReq) (*entities.GetUserRefreshTokensRes, error)
}
