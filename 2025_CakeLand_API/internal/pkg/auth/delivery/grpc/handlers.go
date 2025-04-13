package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/auth"
	gen "2025_CakeLand_API/internal/pkg/auth/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/auth/dto"
	"2025_CakeLand_API/internal/pkg/utils"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type GrpcAuthHandler struct {
	gen.UnimplementedAuthServer

	log        *slog.Logger
	validator  *utils.Validator
	usecase    auth.IAuthUsecase
	mdProvider *md.MetadataProvider
}

func NewGrpcAuthHandler(
	log *slog.Logger,
	validator *utils.Validator,
	usecase auth.IAuthUsecase,
	mdProvider *md.MetadataProvider,
) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		log:        log,
		validator:  validator,
		usecase:    usecase,
		mdProvider: mdProvider,
	}
}

func (h *GrpcAuthHandler) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	// Получение метаданных
	fingerprint, err := h.mdProvider.GetValue(ctx, domains.KeyFingerprint)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyFingerprint),
		)
	}

	// Валидация
	if err = h.validator.ValidateEmail(req.Email); err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "invalid email format")
	} else if err = h.validator.ValidatePassword(req.Password); err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "invalid password format")
	}

	// Сохраняем пользователя в бд
	res, err := h.usecase.Register(ctx, dto.RegisterReq{
		Email:       req.Email,
		Password:    req.Password,
		Fingerprint: fingerprint,
	})
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to register user")
	}

	return &gen.RegisterResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn.Unix(),
	}, nil
}

func (h *GrpcAuthHandler) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	// Получение метаданных
	fingerprint, err := h.mdProvider.GetValue(ctx, domains.KeyFingerprint)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyFingerprint),
		)
	}

	// Валидация
	if err = h.validator.ValidateEmail(req.Email); err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "invalid email format")
	} else if err = h.validator.ValidatePassword(req.Password); err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "invalid password format")
	}

	// Сохраняем пользователя в бд
	res, loginErr := h.usecase.Login(ctx, dto.LoginReq{
		Email:       req.Email,
		Password:    req.Password,
		Fingerprint: fingerprint,
	})
	if loginErr != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, loginErr, "failed to login")
	}

	return &gen.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn.Unix(),
	}, nil
}

func (h *GrpcAuthHandler) Logout(ctx context.Context, _ *emptypb.Empty) (*gen.LogoutResponse, error) {
	// Получение метаданных
	values, err := h.mdProvider.GetValues(ctx, domains.KeyFingerprint, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s or %s", domains.KeyAuthorization, domains.KeyFingerprint),
		)
	}

	fingerprint := values[domains.KeyFingerprint]
	refreshToken := values[domains.KeyAuthorization]

	// Валидация
	if refreshToken == "" {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrNoMetadata, "refreshToken is empty")
	} else if fingerprint == "" {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrNoMetadata, "fingerprint is empty")
	}

	// Бизнес логика
	res, err := h.usecase.Logout(ctx, dto.LogoutReq{
		Fingerprint:  fingerprint,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to logout")
	}

	return &gen.LogoutResponse{
		Message: res.Message,
	}, nil
}

func (h *GrpcAuthHandler) UpdateAccessToken(ctx context.Context, _ *emptypb.Empty) (*gen.UpdateAccessTokenResponse, error) {
	// Получение метаданных
	values, err := h.mdProvider.GetValues(ctx, domains.KeyFingerprint, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s or %s", domains.KeyAuthorization, domains.KeyFingerprint),
		)
	}

	fingerprint := values[domains.KeyFingerprint]
	refreshToken := values[domains.KeyAuthorization]

	// Валидация
	if refreshToken == "" {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrNoMetadata, "refreshToken is empty")
	} else if fingerprint == "" {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrNoMetadata, "fingerprint is empty")
	}

	// Бизнес логика
	res, err := h.usecase.UpdateAccessToken(ctx, dto.UpdateAccessTokenReq{
		RefreshToken: refreshToken,
		Fingerprint:  fingerprint,
	})
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to update access token")
	}

	return &gen.UpdateAccessTokenResponse{
		AccessToken: res.AccessToken,
		ExpiresIn:   res.ExpiresIn.Unix(),
	}, nil
}
