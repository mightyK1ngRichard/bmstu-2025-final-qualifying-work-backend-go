package models

import (
	gen "2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type PaymentMethod string

const (
	Cash    PaymentMethod = "cash"
	IoMoney PaymentMethod = "io_money"
)

type OrderDB struct {
	ID                uuid.UUID
	TotalPrice        float64
	Mass              float64
	PaymentMethod     PaymentMethod
	FillingID         uuid.UUID
	SellerID          uuid.UUID
	CustomerID        uuid.UUID
	CakeID            uuid.UUID
	DeliveryAddressID uuid.UUID
	DeliveryDate      time.Time
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
