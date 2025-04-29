package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/cake"
	gen "2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/cake/dto"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"

	"github.com/google/uuid"
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

func (h *GrpcCakeHandler) AddCakeColors(ctx context.Context, in *gen.AddCakeColorsReq) (*emptypb.Empty, error) {
	cakeID, err := uuid.Parse(in.CakeID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err), "failed to parse cake id")
	}

	if err = h.usecase.AddCakeColor(ctx, cakeID, in.ColorsHex); err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to add color")
	}

	return &emptypb.Empty{}, nil
}

func (h *GrpcCakeHandler) GetColors(ctx context.Context, _ *emptypb.Empty) (*gen.CakeColorsRes, error) {
	// Бизнес логика
	colors, err := h.usecase.GetColors(ctx)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to get colors")
	}

	// Ответ
	return &gen.CakeColorsRes{
		ColorsHex: colors,
	}, nil
}

func (h *GrpcCakeHandler) Cake(ctx context.Context, in *gen.CakeRequest) (*gen.CakeResponse, error) {
	// Параметры
	cakeID, err := uuid.Parse(in.CakeId)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrInvalidUUIDFormat, "'cake_id' must be a valid UUID")
	}

	// Бизнес логика
	res, err := h.usecase.Cake(ctx, dto.GetCakeReq{
		CakeID: cakeID,
	})
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch cake")
	}

	// Ответ
	return &gen.CakeResponse{
		Cake: res.Cake.ConvertToCakeGRPC(),
	}, nil
}

func (h *GrpcCakeHandler) CreateCake(ctx context.Context, in *gen.CreateCakeRequest) (*gen.CreateCakeResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization),
		)
	}

	// Бизнес логика
	res, err := h.usecase.CreateCake(ctx, dto.NewCreateCakeReq(in, accessToken))
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create cake")
	}

	// Ответ
	return &gen.CreateCakeResponse{
		CakeId: res.CakeID,
	}, nil
}

func (h *GrpcCakeHandler) CreateFilling(ctx context.Context, in *gen.CreateFillingRequest) (*gen.CreateFillingResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization),
		)
	}

	// Бизнес логика
	res, err := h.usecase.CreateFilling(ctx, dto.CreateFillingReq{
		Name:        in.Name,
		ImageData:   in.ImageData,
		Content:     in.Content,
		KgPrice:     in.KgPrice,
		Description: in.Description,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create filling")
	}

	// Ответ
	return &gen.CreateFillingResponse{
		Filling: res.Filling.ConvertToFillingGRPC(),
	}, nil
}

func (h *GrpcCakeHandler) CreateCategory(ctx context.Context, in *gen.CreateCategoryRequest) (*gen.CreateCategoryResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization),
		)
	}

	// Бизнес логика
	res, err := h.usecase.CreateCategory(ctx, &dto.CreateCategoryReq{
		Name:        in.Name,
		ImageData:   in.ImageData,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create category")
	}

	// Ответ
	return &gen.CreateCategoryResponse{
		Category: res.Category.ConvertToCategoryGRPC(),
	}, nil
}

func (h *GrpcCakeHandler) Categories(ctx context.Context, _ *emptypb.Empty) (*gen.CategoriesResponse, error) {
	// Бизнес логика
	categories, err := h.usecase.Categories(ctx)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch categories")
	}

	// Маппинг
	categoriesGRPC := make([]*gen.Category, len(*categories))
	for i, it := range *categories {
		categoriesGRPC[i] = it.ConvertToCategoryGRPC()
	}

	// Ответ
	return &gen.CategoriesResponse{
		Categories: categoriesGRPC,
	}, nil
}

func (h *GrpcCakeHandler) Fillings(ctx context.Context, _ *emptypb.Empty) (*gen.FillingsResponse, error) {
	// Бизнес логика
	fillings, err := h.usecase.Fillings(ctx)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch fillings")
	}

	// Маппинг
	fillingsGRPC := make([]*gen.Filling, len(*fillings))
	for i, it := range *fillings {
		fillingsGRPC[i] = it.ConvertToFillingGRPC()
	}

	// Ответ
	return &gen.FillingsResponse{
		Fillings: fillingsGRPC,
	}, nil
}

func (h *GrpcCakeHandler) Cakes(ctx context.Context, _ *emptypb.Empty) (*gen.CakesResponse, error) {
	// Бизнес логика
	cakes, err := h.usecase.Cakes(ctx)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch cakes")
	}

	// Маппинг
	cakesGRPC := make([]*gen.Cake, len(*cakes))
	for i, it := range *cakes {
		cakesGRPC[i] = it.ConvertToCakeGRPC()
	}

	// Ответ
	return &gen.CakesResponse{
		Cakes: cakesGRPC,
	}, nil
}

func (h *GrpcCakeHandler) GetCategoriesByGenderName(ctx context.Context, in *gen.GetCategoriesByGenderNameReq) (*gen.GetCategoriesByGenderNameRes, error) {
	// Параметры
	catGen, err := models.ConvertToCategoryGenderFromGrpc(in.CategoryGender)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "unspecified category gender")
	}

	// Бизнес логика
	categories, err := h.usecase.CategoryIDsByGenderName(ctx, catGen)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch category genders")
	}

	// Маппинг
	res := make([]*gen.Category, len(categories))
	for i, it := range categories {
		res[i] = it.ConvertToCategoryGRPC()
	}

	// Ответ
	return &gen.GetCategoriesByGenderNameRes{
		Categories: res,
	}, nil
}

func (h *GrpcCakeHandler) CategoryPreviewCakes(ctx context.Context, in *gen.CategoryPreviewCakesReq) (*gen.CategoryPreviewCakesRes, error) {
	// Параметры
	categoryID, err := uuid.Parse(in.CategoryID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "parsing category id")
	}

	// Бизнес логика
	previewCakes, err := h.usecase.CategoryPreviewCakes(ctx, categoryID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch preview cakes")
	}

	// Маппинг
	res := make([]*gen.PreviewCake, len(previewCakes))
	for i, it := range previewCakes {
		res[i] = it.ConvertToGrpcModel()
	}

	// Ответ
	return &gen.CategoryPreviewCakesRes{
		PreviewCakes: res,
	}, nil
}
