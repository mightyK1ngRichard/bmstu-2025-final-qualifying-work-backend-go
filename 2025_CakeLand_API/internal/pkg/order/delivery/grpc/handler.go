package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/order"
	gen "2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"log/slog"
)

type OrderHandler struct {
	gen.UnimplementedOrderServiceServer

	log        *slog.Logger
	usecase    order.IOrderUsecase
	mdProvider *md.MetadataProvider
}

func NewOrderHandler(
	log *slog.Logger,
	usecase order.IOrderUsecase,
	mdProvider *md.MetadataProvider,
) *OrderHandler {
	return &OrderHandler{
		log:        log,
		usecase:    usecase,
		mdProvider: mdProvider,
	}
}

func (h *OrderHandler) MakeOrder(ctx context.Context, in *gen.MakeOrderReq) (*gen.MakeOrderRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Маппим модель
	dbOrder, err := models.Init(in)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to init order")
	}

	// Бизнес логика
	createdOrder, err := h.usecase.MakeOrder(ctx, accessToken, dbOrder)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to make order")
	}

	// TODO: Надо сделать уведомление

	// Ответ
	return &gen.MakeOrderRes{
		OrderID: createdOrder.ID.String(),
	}, nil
}

func (h *OrderHandler) getAccessToken(ctx context.Context) (string, error) {
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return "", errs.ConvertToGrpcError(ctx, h.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	return accessToken, nil
}
