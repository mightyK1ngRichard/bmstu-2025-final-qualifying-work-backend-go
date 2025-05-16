package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/cake"
	gen "2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/notification/delivery/grpc/generated"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"

	"github.com/google/uuid"
)

type GrpcCakeHandler struct {
	gen.UnimplementedCakeServiceServer

	log        *slog.Logger
	usecase    cake.ICakeUsecase
	mdProvider *md.MetadataProvider
	nc         generated.NotificationServiceClient
}

func NewCakeHandler(
	logger *slog.Logger,
	uc cake.ICakeUsecase,
	mdProvider *md.MetadataProvider,
	nc generated.NotificationServiceClient,
) *GrpcCakeHandler {
	return &GrpcCakeHandler{
		log:        logger,
		usecase:    uc,
		mdProvider: mdProvider,
		nc:         nc,
	}
}

func (h *GrpcCakeHandler) GetAllCakesWithAllStatuses(ctx context.Context, _ *emptypb.Empty) (*gen.GetAllCakesWithAllStatusesRes, error) {
	// Бизнес логика
	cakes, err := h.usecase.GetCakesPreview(ctx, true)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch cakes")
	}

	// Маппинг
	cakesGRPC := make([]*gen.PreviewCake, len(cakes))
	for i, it := range cakes {
		cakesGRPC[i] = it.ConvertToGrpcModel()
	}

	// Ответ
	return &gen.GetAllCakesWithAllStatusesRes{
		Cakes: cakesGRPC,
	}, nil
}

func (h *GrpcCakeHandler) GetUserCakes(ctx context.Context, in *gen.GetUserCakesReq) (*gen.GetUserCakesRes, error) {
	// Бизнес логика
	cakes, err := h.usecase.GetUserCakes(ctx, in.UserID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch cakes")
	}

	// Маппинг
	cakesGRPC := make([]*gen.PreviewCake, len(cakes))
	for i, it := range cakes {
		cakesGRPC[i] = it.ConvertToGrpcModel()
	}

	// Ответ
	return &gen.GetUserCakesRes{
		Cakes: cakesGRPC,
	}, nil
}

func (h *GrpcCakeHandler) SetCakeVisibility(ctx context.Context, in *gen.SetCakeVisibilityReq) (*emptypb.Empty, error) {
	// Парсим CakeID
	cakeID, err := uuid.Parse(in.CakeID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err), "invalid cake id")
	}

	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization),
		)
	}

	// Бизнес логика
	cakeStatus, err := models.FromProtoCakeStatus(in.Status)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to parse cake status")
	}

	ownerID, cakeName, err := h.usecase.SetCakeVisibility(ctx, accessToken, cakeID, cakeStatus)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "unable to set cake visibility. you must be the owner of the cake to change its visibility")
	}

	// Отправка уведомления
	go func(ownerID, cakeName string, cakeStatus models.CakeStatus) {
		h.sendOrderStatusUpdatedNotification(ctx, cakeName, ownerID, cakeStatus)
	}(ownerID, cakeName, cakeStatus)

	// Ответ
	return &emptypb.Empty{}, nil
}

func (h *GrpcCakeHandler) AddCakeColors(ctx context.Context, in *gen.AddCakeColorsReq) (*emptypb.Empty, error) {
	// Получаем CakeID
	cakeID, err := uuid.Parse(in.CakeID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err), "failed to parse cake id")
	}

	// Бизнес логика
	if err = h.usecase.AddCakeColor(ctx, cakeID, in.ColorsHex); err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to add color")
	}

	// Ответ
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
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization),
		)
	}

	// Бизнес логика
	res, err := h.usecase.Cake(ctx, dto.GetCakeReq{
		CakeID:      cakeID,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch cake")
	}

	// Ответ
	return &gen.CakeResponse{
		Cake:                 res.Cake.ConvertToCakeGRPC(),
		UserCanWriteFeedback: res.CanWriteFeedbacks,
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
	var categoryGenders []models.CategoryGender
	for _, tag := range in.GenderTags {
		gender, err2 := models.ConvertToCategoryGenderFromGrpc(tag)
		if err2 != nil {
			continue
		}
		categoryGenders = append(categoryGenders, gender)
	}
	res, err := h.usecase.CreateCategory(ctx, &dto.CreateCategoryReq{
		Name:            in.Name,
		ImageData:       in.ImageData,
		AccessToken:     accessToken,
		CategoryGenders: categoryGenders,
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
	cakes, err := h.usecase.GetCakesPreview(ctx, false)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch cakes")
	}

	// Маппинг
	cakesGRPC := make([]*gen.PreviewCake, len(cakes))
	for i, it := range cakes {
		cakesGRPC[i] = it.ConvertToGrpcModel()
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
		return nil, errs.ConvertToGrpcError(ctx, h.log, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err), "parsing category id")
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

func (h *GrpcCakeHandler) Add3DModel(ctx context.Context, in *gen.Add3DModelReq) (*gen.Add3DModelRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Бизнес-логика
	modelURL, err := h.usecase.Add3DModel(ctx, accessToken, in)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to add model")
	}

	// Ответ
	return &gen.Add3DModelRes{
		Model3DURL: modelURL,
	}, nil
}

func (h *GrpcCakeHandler) getAccessToken(ctx context.Context) (string, error) {
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return "", errs.ConvertToGrpcError(ctx, h.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	return accessToken, nil
}

// Helpers

func (h *GrpcCakeHandler) sendOrderStatusUpdatedNotification(ctx context.Context, cakeName string, userID string, status models.CakeStatus) {
	// Извлекаем метаданные из родительского контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		h.log.Error("не удалось получить метаданные из контекста")
		return
	}

	// Создаём новый контекст с метаданными из родительского контекста
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	// Формируем заголовок и сообщение
	title := fmt.Sprintf("Статус торта изменён")
	message := fmt.Sprintf("Статус торта \"%s\" изменён на \"%s\".", cakeName, statusToText(status))

	// Формируем запрос
	req := &generated.CreateNotificationRequest{
		Title:       title,
		Message:     message,
		RecipientID: userID,
		Kind:        generated.NotificationKind_ORDER_UPDATE,
	}

	_, err := h.nc.CreateNotification(newCtx, req)
	if err != nil {
		h.log.Error("не удалось отправить уведомление об изменении статуса", "error", err)
	}
}

func statusToText(status models.CakeStatus) string {
	switch status {
	case models.CakeStatusPending:
		return "на рассмотрении"
	case models.CakeStatusApproved:
		return "одобрен"
	case models.CakeStatusRejected:
		return "отклонён"
	case models.CakeStatusHidden:
		return "скрыт"
	default:
		return string(status)
	}
}
