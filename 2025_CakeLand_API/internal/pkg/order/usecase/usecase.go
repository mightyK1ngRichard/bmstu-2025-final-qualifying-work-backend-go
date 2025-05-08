package usecase

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/order"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"context"
	"github.com/google/uuid"
	"time"
)

type OrderUsecase struct {
	tokenator *jwt.Tokenator
	repo      order.IOrderRepository
}

func NewOrderUsecase(
	tokenator *jwt.Tokenator,
	repo order.IOrderRepository,
) *OrderUsecase {
	return &OrderUsecase{
		tokenator: tokenator,
		repo:      repo,
	}
}

func (u *OrderUsecase) MakeOrder(ctx context.Context, accessToken string, dbOrder models.OrderDB) (*models.OrderDB, error) {
	// Достаём UserID
	userID, err := u.getUserUUID(accessToken)
	if err != nil {
		return nil, err
	}

	// Сэтим оставшиеся данные
	dbOrder.ID = uuid.New()
	dbOrder.CustomerID = userID

	// Получение актуальной информации
	cake, err := u.repo.CakeInfo(ctx, dbOrder.CakeID)
	if err != nil {
		return nil, err
	}

	// Получаем актуальную цену торта
	kgPrice := cake.KgPrice
	if cake.DiscountKgPrice.Valid && cake.DiscountEndTime.Valid {
		if cake.DiscountEndTime.Time.After(time.Now()) {
			// Скидка активна
			kgPrice = cake.DiscountKgPrice.Float64
		}
	}

	// Сравниваем итоговую цену
	totalPrice := kgPrice * (dbOrder.Mass / 1000)
	if totalPrice != dbOrder.TotalPrice {
		return nil, errs.ErrTotalPriceIncorrect
	}

	// Запрос в БД
	if err = u.repo.CreateOrder(ctx, dbOrder); err != nil {
		return nil, err
	}

	// Ответ
	return &dbOrder, nil
}

func (u *OrderUsecase) getUserUUID(accessToken string) (uuid.UUID, error) {
	// Достаём UserID
	userIDStr, err := u.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
