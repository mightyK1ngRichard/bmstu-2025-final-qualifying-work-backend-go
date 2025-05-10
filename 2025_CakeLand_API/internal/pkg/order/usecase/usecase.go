package usecase

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/order"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"context"
	"github.com/google/uuid"
	"sync"
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

func (u *OrderUsecase) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Получаем все заказы
	dbOrders, err := u.repo.GetAllOrders(ctx)
	if err != nil {
		return nil, err
	}

	// Получаем начинку и адрес доставки
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	errChan := make(chan error, 1)
	orders := make([]models.Order, len(dbOrders))

	for i, dbOrder := range dbOrders {
		orders[i] = models.MapOrderFromDB(dbOrder)
		iCopy := i
		dbOrderCopy := dbOrder
		wg.Add(2)

		// Получаем данные по адресу
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			address, err2 := u.repo.AddressByID(ctx, dbOrderCopy.DeliveryAddressID)
			if err2 != nil {
				trySendError(err2, errChan, cancel)
				return
			}

			mu.Lock()
			orders[iCopy].DeliveryAddress = *address
			mu.Unlock()
		}()

		// Получаем данные по начинке
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			filling, err2 := u.repo.FillingByID(ctx, dbOrderCopy.FillingID)
			if err2 != nil {
				trySendError(err2, errChan, cancel)
				return
			}

			mu.Lock()
			orders[iCopy].Filling = *filling
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, err
	}

	return orders, nil
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, accessToken string, status models.OrderStatus, orderID string) (string, string, error) {
	// Проверка истёк ли токен
	expired, err := u.tokenator.IsTokenExpired(accessToken, false)
	if err != nil {
		return "", "", err
	} else if expired {
		return "", "", errs.ErrTokenIsExpired
	}

	// Запись в БД
	return u.repo.UpdateOrderStatus(ctx, status, orderID)
}

func (u *OrderUsecase) OrderByID(ctx context.Context, accessToken, orderID string) (*models.Order, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userID, err := u.getUserUUID(accessToken)
	if err != nil {
		return nil, err
	}

	// Получаем заказ из БД
	dbOrder, err := u.repo.OrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if dbOrder.SellerID != userID && dbOrder.CustomerID != userID {
		return nil, errs.ErrForbidden
	}

	orderModel := models.MapOrderFromDB(*dbOrder)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	errChan := make(chan error, 1)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		address, err2 := u.repo.AddressByID(ctx, dbOrder.DeliveryAddressID)
		if err2 != nil {
			trySendError(err2, errChan, cancel)
			return
		}

		mu.Lock()
		orderModel.DeliveryAddress = *address
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		filling, err2 := u.repo.FillingByID(ctx, dbOrder.FillingID)
		if err2 != nil {
			trySendError(err2, errChan, cancel)
			return
		}
		mu.Lock()
		orderModel.Filling = *filling
		mu.Unlock()
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, err
	}

	return &orderModel, nil
}

func (u *OrderUsecase) Orders(ctx context.Context, accessToken string) ([]models.Order, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Достаём UserID
	userID, err := u.getUserUUID(accessToken)
	if err != nil {
		return nil, err
	}

	// Получаем заказы из БД
	dbOrders, err := u.repo.UserOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	errChan := make(chan error, 1)
	orders := make([]models.Order, len(dbOrders))

	for i, dbOrder := range dbOrders {
		orders[i] = models.MapOrderFromDB(dbOrder)
		iCopy := i
		dbOrderCopy := dbOrder
		wg.Add(2)

		// Получаем данные по адресу
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			address, err2 := u.repo.AddressByID(ctx, dbOrderCopy.DeliveryAddressID)
			if err2 != nil {
				trySendError(err2, errChan, cancel)
				return
			}

			mu.Lock()
			orders[iCopy].DeliveryAddress = *address
			mu.Unlock()
		}()

		// Получаем данные по начинке
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			filling, err2 := u.repo.FillingByID(ctx, dbOrderCopy.FillingID)
			if err2 != nil {
				trySendError(err2, errChan, cancel)
				return
			}

			mu.Lock()
			orders[iCopy].Filling = *filling
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, err
	}

	return orders, nil
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

func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}
