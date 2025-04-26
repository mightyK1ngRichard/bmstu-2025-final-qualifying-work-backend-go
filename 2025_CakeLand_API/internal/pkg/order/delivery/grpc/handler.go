package handler

import (
	"2025_CakeLand_API/internal/pkg/order"
	gen "2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	gen.UnimplementedOrderServiceServer

	usecase order.IOrderUsecase
}

func NewOrderHandler(usecase order.IOrderUsecase) *OrderHandler {
	return &OrderHandler{
		usecase: usecase,
	}
}

func (h *OrderHandler) MakeOrder(ctx context.Context, in *gen.MakeOrderReq) (*gen.MakeOrderRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeOrder not implemented")
}
