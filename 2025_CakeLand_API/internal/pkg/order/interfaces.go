package order

import (
	"2025_CakeLand_API/internal/models"
	"context"
	"github.com/google/uuid"
)

type IOrderUsecase interface {
	MakeOrder(context.Context, string, models.OrderDB) (*models.OrderDB, error)
}

type IOrderRepository interface {
	CreateOrder(context.Context, models.OrderDB) error
	CakeInfo(context.Context, uuid.UUID) (models.Cake, error)
}
