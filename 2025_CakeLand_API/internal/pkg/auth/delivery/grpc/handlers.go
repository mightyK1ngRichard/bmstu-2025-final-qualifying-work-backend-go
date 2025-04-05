package handler

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/auth"
	gen "2025_CakeLand_API/internal/pkg/auth/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/auth/dto"
	"2025_CakeLand_API/internal/pkg/utils"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcAuthHandler struct {
	gen.UnimplementedAuthServer

	validator  *utils.Validator
	usecase    auth.IAuthUsecase
	mdProvider *md.MetadataProvider
}

func NewGrpcAuthHandler(
	validator *utils.Validator,
	usecase auth.IAuthUsecase,
	mdProvider *md.MetadataProvider,
) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		validator:  validator,
		usecase:    usecase,
		mdProvider: mdProvider,
	}
}

func (h *GrpcAuthHandler) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	// Получение метаданных
	fingerprint, err := h.mdProvider.GetValue(ctx, md.KeyFingerprint)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "fingerprint отсутствует в метаданных")
	}

	// Валидация
	if err = h.validator.ValidateEmail(req.Email); err != nil {
		return nil, err
	} else if err = h.validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Сохраняем пользователя в бд
	res, err := h.usecase.Register(ctx, dto.RegisterReq{
		Email:       req.Email,
		Password:    req.Password,
		Fingerprint: fingerprint,
	})
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf(`%v`, err))
		}
		return nil, err
	}

	return &gen.RegisterResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn.Unix(),
	}, nil
}

func (h *GrpcAuthHandler) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	// Получение метаданных
	fingerprint, err := h.mdProvider.GetValue(ctx, md.KeyFingerprint)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "fingerprint отсутствует в метаданных")
	}

	// Валидация
	if err = h.validator.ValidateEmail(req.Email); err != nil {
		return nil, err
	} else if err = h.validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Сохраняем пользователя в бд
	res, loginErr := h.usecase.Login(ctx, dto.LoginReq{
		Email:       req.Email,
		Password:    req.Password,
		Fingerprint: fingerprint,
	})
	if loginErr != nil {
		// Преобразование ошибки в формат gRPC
		if errors.Is(loginErr, models.ErrUserNotFound) || errors.Is(loginErr, models.ErrInvalidPassword) {
			return nil, status.Error(codes.NotFound, "неверный логин или пароль")
		}
		return nil, status.Error(codes.Internal, "внутренняя ошибка сервера")
	}

	return &gen.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn.Unix(),
	}, nil
}

func (h *GrpcAuthHandler) Logout(ctx context.Context, _ *emptypb.Empty) (*gen.LogoutResponse, error) {
	// Получение метаданных
	values, err := h.mdProvider.GetValues(ctx, md.KeyFingerprint, md.KeyAuthorization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "fingerprint отсутствует в метаданных")
	}

	fingerprint, fOk := values[md.KeyFingerprint]
	refreshToken, rOk := values[md.KeyAuthorization]

	// Валидация
	if refreshToken == "" || fingerprint == "" || !fOk || !rOk {
		return nil, status.Error(codes.InvalidArgument, "refreshToken обязателен")
	}

	// Бизнес логика
	res, err := h.usecase.Logout(ctx, dto.LogoutReq{
		Fingerprint:  fingerprint,
		RefreshToken: refreshToken,
	})
	if err != nil {
		if errors.Is(err, models.ErrNoToken) {
			return nil, status.Error(codes.InvalidArgument, "неверный refresh токен")
		}
		return nil, err
	}

	return &gen.LogoutResponse{
		Message: res.Message,
	}, nil
}

func (h *GrpcAuthHandler) UpdateAccessToken(ctx context.Context, _ *emptypb.Empty) (*gen.UpdateAccessTokenResponse, error) {
	// Получение метаданных
	values, err := h.mdProvider.GetValues(ctx, md.KeyFingerprint, md.KeyAuthorization)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "fingerprint отсутствует в метаданных")
	}

	fingerprint := values[md.KeyFingerprint]
	refreshToken := values[md.KeyAuthorization]

	// Валидация
	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refreshToken обязателен")
	}

	// Бизнес логика
	res, err := h.usecase.UpdateAccessToken(ctx, dto.UpdateAccessTokenReq{
		RefreshToken: refreshToken,
		Fingerprint:  fingerprint,
	})
	if err != nil {
		return nil, err
	}

	return &gen.UpdateAccessTokenResponse{
		AccessToken: res.AccessToken,
		ExpiresIn:   res.ExpiresIn.Unix(),
	}, nil
}
