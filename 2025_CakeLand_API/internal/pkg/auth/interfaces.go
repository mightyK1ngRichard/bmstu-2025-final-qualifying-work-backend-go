package auth

import (
	"2025_CakeLand_API/internal/pkg/auth/dto"
	"context"
)

// mockgen -source=internal/pkg/auth/interfaces.go -destination=internal/pkg/auth/mocks/mock_auth.go -package=mocks

type IAuthUsecase interface {
	Register(context.Context, dto.RegisterReq) (*dto.RegisterRes, error)
	Login(context.Context, dto.LoginReq) (*dto.LoginRes, error)
	Logout(context.Context, dto.LogoutReq) (*dto.LogoutRes, error)
	UpdateAccessToken(context.Context, dto.UpdateAccessTokenReq) (*dto.UpdateAccessTokenRes, error)
}

type IAuthRepository interface {
	CreateUser(context.Context, dto.CreateUserReq) error
	GetUserByEmail(context.Context, dto.GetUserByEmailReq) (*dto.GetUserByEmailRes, error)
	UpdateUserRefreshTokens(context.Context, dto.UpdateUserRefreshTokensReq) error
	GetUserRefreshTokens(context.Context, dto.GetUserRefreshTokensReq) (*dto.GetUserRefreshTokensRes, error)
}
