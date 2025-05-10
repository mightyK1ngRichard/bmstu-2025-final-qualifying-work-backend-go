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
	queryGetAllOrders = `
		SELECT id,
			   total_price,
			   delivery_address_id,
			   mass,
			   filling_id,
			   delivery_date,
			   customer_id,
			   seller_id,
			   cake_id,
			   payment_method,
			   status
		FROM "order"
	`
	queryUpdateOrderStatus = `UPDATE "order" SET status = $1 WHERE id = $2 RETURNING customer_id, seller_id;`
	queryCreateOrder       = `
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
	queryUserOrders = `
		SELECT id,
			   total_price,
			   delivery_address_id,
			   mass,
			   filling_id,
			   delivery_date,
			   customer_id,
			   seller_id,
			   payment_method,
			   cake_id,
			   status
		FROM "order"
		WHERE customer_id = $1
		ORDER BY delivery_date DESC
	`
	queryOrderByID = `
		SELECT id,
		   total_price,
		   delivery_address_id,
		   mass,
		   filling_id,
		   delivery_date,
		   customer_id,
		   seller_id,
		   payment_method,
		   cake_id,
		   status
		FROM "order"
		WHERE id = $1
		LIMIT 1
	`
	queryAddressByID = `
		SELECT id,
			   user_id,
			   latitude,
			   longitude,
			   formatted_address,
			   entrance,
			   floor,
			   apartment,
			   comment
		FROM address
		WHERE id = $1
	`
	queryGetFillingByID = `SELECT id, name, image_url, content, kg_price, description FROM filling WHERE id = $1`
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) OrderByID(ctx context.Context, orderID string) (*models.OrderDB, error) {
	const methodName = "[OrderRepo.OrderByID]"

	row := r.db.QueryRowContext(ctx, queryOrderByID, orderID)

	var o models.OrderDB
	if err := row.Scan(
		&o.ID,
		&o.TotalPrice,
		&o.DeliveryAddressID,
		&o.Mass,
		&o.FillingID,
		&o.DeliveryDate,
		&o.CustomerID,
		&o.SellerID,
		&o.PaymentMethod,
		&o.CakeID,
		&o.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.WrapDBError(methodName, err)
	}

	return &o, nil
}

func (r *OrderRepo) GetAllOrders(ctx context.Context) ([]models.OrderDB, error) {
	const methodName = "[OrderRepo.GetAllOrders]"

	rows, err := r.db.QueryContext(ctx, queryGetAllOrders)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}
	defer rows.Close()

	var orders []models.OrderDB
	for rows.Next() {
		var o models.OrderDB
		if err = rows.Scan(
			&o.ID,
			&o.TotalPrice,
			&o.DeliveryAddressID,
			&o.Mass,
			&o.FillingID,
			&o.DeliveryDate,
			&o.CustomerID,
			&o.SellerID,
			&o.CakeID,
			&o.PaymentMethod,
			&o.Status,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return orders, nil
}

func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, status models.OrderStatus, orderID string) (string, string, error) {
	const methodName = "[OrderRepo.UpdateOrderStatus]"

	var customerID string
	var sellerID string
	if err := r.db.QueryRowContext(ctx, queryUpdateOrderStatus, status, orderID).Scan(&customerID, &sellerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errs.ErrNotFound
		}
		return "", "", errs.WrapDBError(methodName, err)
	}

	return customerID, sellerID, nil
}

func (r *OrderRepo) FillingByID(ctx context.Context, fillingID uuid.UUID) (*models.Filling, error) {
	const methodName = "[OrderRepo.FillingByID]"

	var filling models.Filling
	if err := r.db.QueryRowContext(ctx, queryGetFillingByID, fillingID).Scan(
		&filling.ID,
		&filling.Name,
		&filling.ImageURL,
		&filling.Content,
		&filling.KgPrice,
		&filling.Description,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, errs.WrapDBError(methodName, err)
	}

	return &filling, nil
}

func (r *OrderRepo) AddressByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	const methodName = "[OrderRepo.AddressByID]"

	row := r.db.QueryRowContext(ctx, queryAddressByID, id)

	var addr models.Address
	if err := row.Scan(
		&addr.ID,
		&addr.UserID,
		&addr.Latitude,
		&addr.Longitude,
		&addr.FormattedAddress,
		&addr.Entrance,
		&addr.Floor,
		&addr.Apartment,
		&addr.Comment,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.WrapDBError(methodName, err)
	}

	return &addr, nil
}

func (r *OrderRepo) UserOrders(ctx context.Context, userID uuid.UUID) ([]models.OrderDB, error) {
	const methodName = "[OrderRepo.UserOrders]"

	rows, err := r.db.QueryContext(ctx, queryUserOrders, userID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}
	defer rows.Close()

	var orders []models.OrderDB
	for rows.Next() {
		var o models.OrderDB
		if err = rows.Scan(
			&o.ID,
			&o.TotalPrice,
			&o.DeliveryAddressID,
			&o.Mass,
			&o.FillingID,
			&o.DeliveryDate,
			&o.CustomerID,
			&o.SellerID,
			&o.PaymentMethod,
			&o.CakeID,
			&o.Status,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return orders, nil
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
