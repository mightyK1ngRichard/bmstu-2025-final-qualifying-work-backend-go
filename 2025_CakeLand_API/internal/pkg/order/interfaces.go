package order

import (
	"2025_CakeLand_API/internal/models"
	"context"
	"github.com/google/uuid"
)

type IOrderUsecase interface {
	MakeOrder(context.Context, string, models.OrderDB) (*models.OrderDB, error)
	Orders(context.Context, string) ([]models.Order, error)
	UpdateOrderStatus(context.Context, models.OrderStatus, string) (string, string, error)
	GetAllOrders(context.Context) ([]models.Order, error)
	OrderByID(context.Context, string, string) (*models.Order, error)
}

type IOrderRepository interface {
	CreateOrder(context.Context, models.OrderDB) error
	CakeInfo(context.Context, uuid.UUID) (models.Cake, error)
	UserOrders(context.Context, uuid.UUID) ([]models.OrderDB, error)
	AddressByID(context.Context, uuid.UUID) (*models.Address, error)
	FillingByID(context.Context, uuid.UUID) (*models.Filling, error)
	UpdateOrderStatus(context.Context, models.OrderStatus, string) (string, string, error)
	GetAllOrders(context.Context) ([]models.OrderDB, error)
	OrderByID(context.Context, string) (*models.OrderDB, error)
}
