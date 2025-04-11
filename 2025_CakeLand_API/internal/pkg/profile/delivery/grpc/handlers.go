package handler

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/profile"
	gen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	accessToken, err := h.mdProvider.GetValue(ctx, md.KeyAuthorization)
	if err != nil {
		h.log.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userInfo, err := h.usecase.UserInfo(ctx, accessToken)
	if err != nil {
		return nil, models.HandleError(ctx, h.log, err, "failed to get user info")
	}

	return &gen.GetUserInfoRes{
		UserInfo: userInfo.ConvertToGrpcModel(),
	}, nil
}
