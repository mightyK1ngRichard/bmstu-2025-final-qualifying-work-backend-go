package models

import gen "2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) String() string {
	return string(s)
}

// Proto -> GO
func InitFromProtoOrderStatus(status gen.OrderStatus) OrderStatus {
	switch status {
	case gen.OrderStatus_PENDING:
		return OrderStatusPending
	case gen.OrderStatus_SHIPPED:
		return OrderStatusShipped
	case gen.OrderStatus_DELIVERED:
		return OrderStatusDelivered
	case gen.OrderStatus_CANCELLED:
		return OrderStatusCancelled
	default:
		return OrderStatusPending
	}
}

// Go -> Proto
func toProtoOrderStatus(status OrderStatus) gen.OrderStatus {
	switch status {
	case OrderStatusPending:
		return gen.OrderStatus_PENDING
	case OrderStatusShipped:
		return gen.OrderStatus_SHIPPED
	case OrderStatusDelivered:
		return gen.OrderStatus_DELIVERED
	case OrderStatusCancelled:
		return gen.OrderStatus_CANCELLED
	default:
		return gen.OrderStatus_PENDING
	}
}
