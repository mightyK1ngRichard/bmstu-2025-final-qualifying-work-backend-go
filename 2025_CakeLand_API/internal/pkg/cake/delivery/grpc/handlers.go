package handler

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/cake"
	gen "2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	en "2025_CakeLand_API/internal/pkg/cake/entities"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcCakeHandler struct {
	gen.UnimplementedCakeServiceServer

	log        *slog.Logger
	usecase    cake.ICakeUsecase
	mdProvider *md.MetadataProvider
}

func NewCakeHandler(
	logger *slog.Logger,
	uc cake.ICakeUsecase,
	mdProvider *md.MetadataProvider,
) *GrpcCakeHandler {
	return &GrpcCakeHandler{
		log:        logger,
		usecase:    uc,
		mdProvider: mdProvider,
	}
}

func (h *GrpcCakeHandler) Cake(ctx context.Context, in *gen.CakeRequest) (*gen.CakeResponse, error) {
	cakeID, err := uuid.Parse(in.CakeId)
	if err != nil {
		h.log.Error("Ошибка парсинга CakeID", "CakeID", in.CakeId, "error", err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Некорректный формат CakeID: %s", in.CakeId))
	}

	res, err := h.usecase.Cake(ctx, en.GetCakeReq{
		CakeID: cakeID,
	})
	if err = models.HandleError(err); err != nil {
		return nil, err
	}

	// Формируем CakeResponse
	return &gen.CakeResponse{
		Cake: res.Cake.ConvertToCakeGRPC(),
	}, nil
}

func (h *GrpcCakeHandler) CreateCake(ctx context.Context, in *gen.CreateCakeRequest) (*gen.CreateCakeResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, md.KeyAuthorization)
	if err != nil {
		h.log.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := h.usecase.CreateCake(ctx, en.NewCreateCakeReq(in, accessToken))
	if err = models.HandleError(err); err != nil {
		return nil, err
	}

	return &gen.CreateCakeResponse{
		CakeId: res.CakeID,
	}, nil
}

func (h *GrpcCakeHandler) CreateFilling(ctx context.Context, in *gen.CreateFillingRequest) (*gen.CreateFillingResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, md.KeyAuthorization)
	if err != nil {
		h.log.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := h.usecase.CreateFilling(ctx, en.CreateFillingReq{
		Name:        in.Name,
		ImageData:   in.ImageData,
		Content:     in.Content,
		KgPrice:     in.KgPrice,
		Description: in.Description,
		AccessToken: accessToken,
	})
	if err = models.HandleError(err); err != nil {
		return nil, err
	}

	return &gen.CreateFillingResponse{
		Filling: res.Filling.ConvertToFillingGRPC(),
	}, nil
}

func (h *GrpcCakeHandler) CreateCategory(ctx context.Context, in *gen.CreateCategoryRequest) (*gen.CreateCategoryResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, md.KeyAuthorization)
	if err != nil {
		h.log.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := h.usecase.CreateCategory(ctx, &en.CreateCategoryReq{
		Name:        in.Name,
		ImageData:   in.ImageData,
		AccessToken: accessToken,
	})
	if err = models.HandleError(err); err != nil {
		return nil, err
	}

	return &gen.CreateCategoryResponse{
		Category: res.Category.ConvertToCategoryGRPC(),
	}, nil
}

func (h *GrpcCakeHandler) Categories(ctx context.Context, _ *emptypb.Empty) (*gen.CategoriesResponse, error) {
	categories, err := h.usecase.Categories(ctx)
	if err = models.HandleError(err); err != nil {
		return nil, err
	}

	categoriesGRPC := make([]*gen.Category, len(*categories))
	for i, it := range *categories {
		categoriesGRPC[i] = it.ConvertToCategoryGRPC()
	}

	return &gen.CategoriesResponse{
		Categories: categoriesGRPC,
	}, nil
}

func (h *GrpcCakeHandler) Fillings(ctx context.Context, _ *emptypb.Empty) (*gen.FillingsResponse, error) {
	fillings, err := h.usecase.Fillings(ctx)
	if err = models.HandleError(err); err != nil {
		return nil, err
	}

	fillingsGRPC := make([]*gen.Filling, len(*fillings))
	for i, it := range *fillings {
		fillingsGRPC[i] = it.ConvertToFillingGRPC()
	}

	return &gen.FillingsResponse{
		Fillings: fillingsGRPC,
	}, nil
}

func (h *GrpcCakeHandler) Cakes(ctx context.Context, _ *emptypb.Empty) (*gen.CakesResponse, error) {
	cakes, err := h.usecase.Cakes(ctx)
	if err = models.HandleError(err); err != nil {
		return nil, err
	}

	cakesGRPC := make([]*gen.Cake, len(*cakes))
	for i, it := range *cakes {
		cakesGRPC[i] = it.ConvertToCakeGRPC()
	}

	return &gen.CakesResponse{
		Cakes: cakesGRPC,
	}, nil
}
