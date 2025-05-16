package models

import gen "2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"

type PaymentMethod string

const (
	Cash    PaymentMethod = "cash"
	IoMoney PaymentMethod = "io_money"
)

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
