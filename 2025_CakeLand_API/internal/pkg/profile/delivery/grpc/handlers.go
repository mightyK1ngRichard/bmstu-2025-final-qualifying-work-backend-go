package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/profile"
	gen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type GrpcProfileHandler struct {
	gen.UnimplementedProfileServiceServer

	log        *slog.Logger
	usecase    profile.IProfileUsecase
	mdProvider *md.MetadataProvider
}

func NewProfileHandler(
	logger *slog.Logger,
	uc profile.IProfileUsecase,
	mdProvider *md.MetadataProvider,
) *GrpcProfileHandler {
	return &GrpcProfileHandler{
		log:        logger,
		usecase:    uc,
		mdProvider: mdProvider,
	}
}

func (h *GrpcProfileHandler) GetUserInfo(ctx context.Context, _ *emptypb.Empty) (*gen.GetUserInfoRes, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	userInfo, err := h.usecase.UserInfo(ctx, accessToken)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch user info")
	}

	return &gen.GetUserInfoRes{
		UserInfo: userInfo.ConvertToGrpcModel(),
	}, nil
}

func (h *GrpcProfileHandler) GetUserInfoByID(ctx context.Context, req *gen.GetUserInfoByIDReq) (*gen.GetUserInfoByIDRes, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err), "invalid user id format")
	}

	userInfo, err := h.usecase.UserInfoByID(ctx, userID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch user info")
	}

	return &gen.GetUserInfoByIDRes{
		User: userInfo.ConvertToGRPCProfile(),
	}, nil
}
