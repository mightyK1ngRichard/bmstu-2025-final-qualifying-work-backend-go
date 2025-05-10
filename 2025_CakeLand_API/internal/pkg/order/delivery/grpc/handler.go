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
	"google.golang.org/protobuf/types/known/emptypb"
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

func (h *OrderHandler) GetAllOrders(ctx context.Context, _ *emptypb.Empty) (*gen.GetAllOrdersRes, error) {
	// Бизнес логика
	orders, err := h.usecase.GetAllOrders(ctx)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to get all orders")
	}

	// Ответ
	res := make([]*gen.Order, len(orders))
	for i, orderItem := range orders {
		res[i] = orderItem.ToProto()
	}

	return &gen.GetAllOrdersRes{
		Orders: res,
	}, nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, in *gen.UpdateOrderStateReq) (*emptypb.Empty, error) {
	// Получаем токен из метаданных
	_, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Получаем статус
	updatedStatus := models.InitFromProtoOrderStatus(in.UpdatedStatus)

	// Бизнес логика
	customerID, sellerID, err := h.usecase.UpdateOrderStatus(ctx, updatedStatus, in.OrderID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to update order")
	}

	// Уведомление
	go func() {
		h.sendOrderStatusUpdatedNotification(ctx, customerID, in.OrderID, updatedStatus)
	}()
	go func() {
		h.sendOrderStatusUpdatedNotification(ctx, sellerID, in.OrderID, updatedStatus)
	}()

	// Ответ
	return &emptypb.Empty{}, nil
}

func (h *OrderHandler) Orders(ctx context.Context, _ *emptypb.Empty) (*gen.OrdersRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Бизнес логика
	orders, err := h.usecase.Orders(ctx, accessToken)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch orders")
	}

	// Ответ
	res := make([]*gen.Order, len(orders))
	for i, orderItem := range orders {
		res[i] = orderItem.ToProto()
	}

	return &gen.OrdersRes{
		Orders: res,
	}, nil
}

func (h *OrderHandler) OrderByID(ctx context.Context, in *gen.OrderByIDReq) (*gen.OrderByIDRes, error) {
	// Получаем токен из метаданных
	accessToken, convertedErr := h.getAccessToken(ctx)
	if convertedErr != nil {
		return nil, convertedErr
	}

	// Бизнес логика
	orderModel, err := h.usecase.OrderByID(ctx, accessToken, in.OrderID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to fetch order by ID")
	}

	// Ответ
	return &gen.OrderByIDRes{
		Order: orderModel.ToProto(),
	}, nil
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
		// Отправка уведомления для продавца
		h.sendOrderNotification(ctx, createdOrder.SellerID.String(), createdOrder.ID.String(),
			"Заказ оформлен", "Ваш заказ успешно сформирован и находится в обработке 🎂")
	}()

	go func() {
		// Отправка уведомления для покупателя
		h.sendOrderNotification(ctx, createdOrder.CustomerID.String(), createdOrder.ID.String(),
			"Торт заказан", "У вас заказали торт! 🎂 Ваш заказ находится в обработке.")
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

func (h *OrderHandler) sendOrderNotification(ctx context.Context, userID, orderID, messageTitle, messageBody string) {
	// Извлекаем метаданные из родительского контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		h.log.Error("не удалось получить метаданные из контекста")
		return
	}

	// Создаём новый контекст с метаданными из родительского контекста
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	req := &generated.CreateNotificationRequest{
		Title:       messageTitle,
		Message:     messageBody,
		OrderID:     &orderID,
		RecipientID: userID,
		Kind:        generated.NotificationKind_ORDER_UPDATE,
	}

	_, err := h.nc.CreateNotification(newCtx, req)
	if err != nil {
		h.log.Error("не удалось отправить уведомление", "error", err)
	}
}

func (h *OrderHandler) sendOrderStatusUpdatedNotification(ctx context.Context, userID, orderID string, status models.OrderStatus) {
	// Извлекаем метаданные из родительского контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		h.log.Error("не удалось получить метаданные из контекста")
		return
	}

	// Создаём новый контекст с метаданными из родительского контекста
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	// Получаем заголовок и сообщение в зависимости от статуса заказа
	title, message := getStatusNotificationText(status)
	req := &generated.CreateNotificationRequest{
		Title:       title,
		Message:     message,
		OrderID:     &orderID,
		RecipientID: userID,
		Kind:        generated.NotificationKind_ORDER_UPDATE,
	}

	_, err := h.nc.CreateNotification(newCtx, req)
	if err != nil {
		h.log.Error("не удалось отправить уведомление об изменении статуса", "error", err)
	}
}

func getStatusNotificationText(status models.OrderStatus) (title, message string) {
	switch status {
	case models.OrderStatusPending:
		return "Ожидание обработки", "Ваш заказ ожидает обработки 🍰"
	case models.OrderStatusShipped:
		return "Заказ в пути", "Ваш заказ уже в пути к вам 🚚"
	case models.OrderStatusDelivered:
		return "Доставка завершена", "Ваш заказ доставлен, приятного аппетита! 🎉"
	case models.OrderStatusCancelled:
		return "Заказ отменён", "Ваш заказ был отменён ❌"
	default:
		return "Обновление заказа", "Статус вашего заказа обновлён"
	}
}
