package usecase

import "context"

type OrderUsecase struct{}

func NewOrderUsecase() *OrderUsecase {
	return &OrderUsecase{}
}

func (u *OrderUsecase) MakeOrder(ctx context.Context) error {
	return nil
}
