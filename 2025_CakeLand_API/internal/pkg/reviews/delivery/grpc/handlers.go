package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/reviews"
	gen "2025_CakeLand_API/internal/pkg/reviews/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/reviews/entities"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

type GrpcReviewsHandler struct {
	gen.UnimplementedReviewServiceServer

	log        *slog.Logger
	usecase    reviews.IReviewsUsecase
	mdProvider *md.MetadataProvider
	tokenator  *jwt.Tokenator
}

func NewReviewsHandler(
	logger *slog.Logger,
	uc reviews.IReviewsUsecase,
	mdProvider *md.MetadataProvider,
	tokenator *jwt.Tokenator,
) *GrpcReviewsHandler {
	return &GrpcReviewsHandler{
		log:        logger,
		usecase:    uc,
		mdProvider: mdProvider,
		tokenator:  tokenator,
	}
}

func (h *GrpcReviewsHandler) AddFeedback(ctx context.Context, in *gen.AddFeedbackRequest) (*gen.AddFeedbackResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	userID, err := h.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, fmt.Sprintf("missing required token: %s", domains.KeyAuthorization))
	}

	// Валидация
	if in.Text == "" {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrInvalidInput, "text is required")
	}
	if !(in.Rating > 0 && in.Rating < 6) {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrInvalidInput, "rating must be between 1 and 5")
	}

	// Бизнес логика
	request, err := entities.NewCreateFeedbackReq(in, userID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create request")
	}

	feedback, err := h.usecase.CreateFeedback(ctx, *request)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create request")
	}

	// Ответ
	return &gen.AddFeedbackResponse{
		Feedback: feedback.ConvertToGRPC(),
	}, nil
}

func (h *GrpcReviewsHandler) ProductFeedbacks(ctx context.Context, in *gen.ProductFeedbacksRequest) (*gen.ProductFeedbacksResponse, error) {
	cakeID, err := uuid.Parse(in.CakeID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, errs.ErrInvalidUUIDFormat, "invalid cake id")
	}

	feedbacks, err := h.usecase.ProductFeedbacks(ctx, cakeID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create request")
	}

	response := make([]*gen.Feedback, len(feedbacks))
	for i, feedback := range feedbacks {
		response[i] = feedback.ConvertToGRPC()
	}

	return &gen.ProductFeedbacksResponse{
		Feedbacks: response,
	}, nil
}
