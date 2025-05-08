package models

import (
	gen "2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type PaymentMethod string
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"

	Cash    PaymentMethod = "cash"
	IoMoney PaymentMethod = "io_money"
)

func (s OrderStatus) String() string {
	return string(s)
}

type Order struct {
	ID              uuid.UUID
	TotalPrice      float64
	DeliveryAddress Address
	Mass            float64
	Filling         Filling
	DeliveryDate    time.Time
	SellerID        uuid.UUID
	PaymentMethod   PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CakeID          uuid.UUID
}

type OrderDB struct {
	ID                uuid.UUID
	TotalPrice        float64
	Mass              float64
	Status            OrderStatus
	PaymentMethod     PaymentMethod
	FillingID         uuid.UUID
	SellerID          uuid.UUID
	CustomerID        uuid.UUID
	CakeID            uuid.UUID
	DeliveryAddressID uuid.UUID
	DeliveryDate      time.Time
}

func MapOrderFromDB(dbOrder OrderDB) Order {
	return Order{
		ID:            dbOrder.ID,
		TotalPrice:    dbOrder.TotalPrice,
		Mass:          dbOrder.Mass,
		DeliveryDate:  dbOrder.DeliveryDate,
		SellerID:      dbOrder.SellerID,
		PaymentMethod: dbOrder.PaymentMethod,
		Status:        dbOrder.Status,
		CakeID:        dbOrder.CakeID,
	}
}

func Init(from *gen.MakeOrderReq) (OrderDB, error) {
	deliveryDate := time.Time{}
	if from.DeliveryDate != nil {
		deliveryDate = from.DeliveryDate.AsTime()
	}

	deliveryAddressID, err := uuid.Parse(from.DeliveryAddressID)
	if err != nil {
		return OrderDB{}, err
	}

	fillingID, err := uuid.Parse(from.FillingID)
	if err != nil {
		return OrderDB{}, err
	}

	sellerID, err := uuid.Parse(from.SellerID)
	if err != nil {
		return OrderDB{}, err
	}

	cakeID, err := uuid.Parse(from.CakeID)
	if err != nil {
		return OrderDB{}, err
	}

	var paymentMethod PaymentMethod
	switch from.PaymentMethod {
	case gen.PaymentMethod_CASH:
		paymentMethod = Cash
	case gen.PaymentMethod_IOMoney:
		paymentMethod = IoMoney
	default:
		return OrderDB{}, fmt.Errorf("unknown payment method: %v", from.PaymentMethod)
	}

	return OrderDB{
		ID:                uuid.New(),
		TotalPrice:        from.TotalPrice,
		Mass:              from.Mass,
		PaymentMethod:     paymentMethod,
		FillingID:         fillingID,
		SellerID:          sellerID,
		CakeID:            cakeID,
		DeliveryAddressID: deliveryAddressID,
		DeliveryDate:      deliveryDate,
	}, nil
}

func (o *Order) ToProto() *gen.Order {
	return &gen.Order{
		Id:              o.ID.String(),
		TotalPrice:      o.TotalPrice,
		DeliveryAddress: o.DeliveryAddress.ConvertToGRPCAddress(),
		Mass:            o.Mass,
		Filling:         o.Filling.ConvertToFillingGRPC(),
		DeliveryDate:    timestamppb.New(o.DeliveryDate),
		SellerID:        o.SellerID.String(),
		CakeID:          o.CakeID.String(),
		PaymentMethod:   toProtoPaymentMethod(o.PaymentMethod),
		Status:          toProtoOrderStatus(o.Status),
		CreatedAt:       timestamppb.New(o.CreatedAt),
		UpdatedAt:       timestamppb.New(o.UpdatedAt),
	}
}

func toProtoPaymentMethod(pm PaymentMethod) gen.PaymentMethod {
	switch pm {
	case Cash:
		return gen.PaymentMethod_CASH
	case IoMoney:
		return gen.PaymentMethod_IOMoney
	default:
		return gen.PaymentMethod_CASH
	}
}

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
