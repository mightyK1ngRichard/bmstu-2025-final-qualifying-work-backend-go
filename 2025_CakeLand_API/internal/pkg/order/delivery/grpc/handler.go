package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/notification/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/order"
	gen "2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

type OrderHandler struct {
	gen.UnimplementedOrderServiceServer

	log        *slog.Logger
	usecase    order.IOrderUsecase
	mdProvider *md.MetadataProvider
	nc         generated.NotificationServiceClient
}

func NewOrderHandler(
	log *slog.Logger,
	usecase order.IOrderUsecase,
	mdProvider *md.MetadataProvider,
	nc generated.NotificationServiceClient,
) *OrderHandler {
	return &OrderHandler{
		log:        log,
		usecase:    usecase,
		mdProvider: mdProvider,
		nc:         nc,
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

	// Отправка уведомления
	go func() {
		h.sendOrderCreatedNotification(ctx, createdOrder.SellerID.String(), createdOrder.CakeID.String())
	}()

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

func (h *OrderHandler) sendOrderCreatedNotification(ctx context.Context, userID, cakeID string) {
	// Извлекаем метаданные из родительского контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		h.log.Error("не удалось получить метаданные из контекста")
		return
	}

	// Создаём новый контекст с метаданными из родительского контекста
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	req := &generated.CreateNotificationRequest{
		Title:       "Заказ оформлен",
		Message:     "Ваш заказ успешно сформирован и находится в обработке 🎂",
		CakeID:      &cakeID, // FIXME: Я хочу отправлять ID заказа
		RecipientID: userID,
		Kind:        generated.NotificationKind_ORDER_UPDATE,
	}

	_, err := h.nc.CreateNotification(newCtx, req)
	if err != nil {
		h.log.Error("не удалось отправить уведомление", "error", err)
	}
}
