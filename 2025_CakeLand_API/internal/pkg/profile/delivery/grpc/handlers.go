package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
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

func (h *GrpcProfileHandler) UpdateUserAddresses(ctx context.Context, req *gen.UpdateUserAddressesReq) (*gen.UpdateUserAddressesRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Бизнес-логика
	address, err := h.usecase.UpdateUserAddresses(ctx, accessToken, req)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to update user address")
	}

	// Ответ
	return &gen.UpdateUserAddressesRes{
		Address: address.ConvertToGRPCAddress(),
	}, nil
}

func (h *GrpcProfileHandler) GetUserAddresses(ctx context.Context, _ *emptypb.Empty) (*gen.GetUserAddressesRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Бизнес-логика
	addresses, err := h.usecase.GetUserAddresses(ctx, accessToken)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to get user addresses")
	}

	// Ответ
	grpcAddresses := make([]*gen.Address, len(addresses))
	for i, address := range addresses {
		grpcAddresses[i] = address.ConvertToGRPCAddress()
	}

	return &gen.GetUserAddressesRes{
		Addresses: grpcAddresses,
	}, nil
}

func (h *GrpcProfileHandler) CreateAddress(ctx context.Context, in *gen.CreateAddressReq) (*gen.CreateAddressRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Бизнес-логика
	address := &models.Address{
		Latitude:         in.Latitude,
		Longitude:        in.Longitude,
		FormattedAddress: in.FormattedAddress,
	}
	createdAddress, err := h.usecase.CreateAddress(ctx, accessToken, address)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create address")
	}

	// Ответ
	return &gen.CreateAddressRes{
		Address: createdAddress.ConvertToGRPCAddress(),
	}, nil
}

func (h *GrpcProfileHandler) GetUserInfo(ctx context.Context, _ *emptypb.Empty) (*gen.GetUserInfoRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Бизнес-логика
	userInfo, err := h.usecase.UserInfo(ctx, accessToken)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch user info")
	}

	// Ответ
	return &gen.GetUserInfoRes{
		UserInfo: userInfo.ConvertToGrpcModel(),
	}, nil
}

func (h *GrpcProfileHandler) GetUserInfoByID(ctx context.Context, req *gen.GetUserInfoByIDReq) (*gen.GetUserInfoByIDRes, error) {
	// Получаем токен из метаданных
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err), "invalid user id format")
	}

	// Бизнес-логика
	userInfo, err := h.usecase.UserInfoByID(ctx, userID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch user info")
	}

	// Ответ
	return &gen.GetUserInfoByIDRes{
		User: userInfo.ConvertToGRPCProfile(),
	}, nil
}

func (h *GrpcProfileHandler) getAccessToken(ctx context.Context) (string, error) {
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return "", errs.ConvertToGrpcError(ctx, h.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	return accessToken, nil
}
