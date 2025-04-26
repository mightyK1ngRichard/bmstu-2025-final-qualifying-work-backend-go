package repo

import (
	"context"
	"database/sql"
)

const (
	queryCreateOrder = ``
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) CreateOrder(ctx context.Context) error {

	return nil
}
