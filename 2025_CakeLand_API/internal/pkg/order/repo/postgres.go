package repo

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

const (
	queryCreateOrder = `
		INSERT INTO "order" (id,
                     total_price,
                     delivery_address_id,
                     mass,
                     filling_id,
                     delivery_date,
                     customer_id,
                     seller_id,
					 payment_method,
                     cake_id)
		VALUES ($1, $2, $3, $4, $5, $6 , $7, $8, $9, $10) 
	`
	queryCakeInfo = `
		SELECT kg_price, mass, discount_kg_price, discount_end_time, is_open_for_sale
		FROM cake
		WHERE id = $1
	`
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, in models.OrderDB) error {
	const methodName = "[OrderRepo.CreateOrder]"

	if _, err := r.db.ExecContext(ctx, queryCreateOrder,
		in.ID,
		in.TotalPrice,
		in.DeliveryAddressID,
		in.Mass,
		in.FillingID,
		in.DeliveryDate,
		in.CustomerID,
		in.SellerID,
		in.PaymentMethod,
		in.CakeID,
	); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *OrderRepo) CakeInfo(ctx context.Context, cakeID uuid.UUID) (models.Cake, error) {
	const methodName = "[OrderRepo.CakeInfo]"

	var cake models.Cake
	if err := r.db.QueryRowContext(ctx, queryCakeInfo, cakeID).Scan(
		&cake.KgPrice,
		&cake.Mass,
		&cake.DiscountKgPrice,
		&cake.DiscountEndTime,
		&cake.IsOpenForSale,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Cake{}, errs.ErrNotFound
		}
		return models.Cake{}, errs.WrapDBError(methodName, err)
	}

	return cake, nil
}
