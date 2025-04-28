package models

import (
	gen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type Address struct {
	ID               uuid.UUID   // Код адреса
	UserID           uuid.UUID   // Код пользователя
	Latitude         float64     // Широта
	Longitude        float64     // Долгота
	FormattedAddress string      // Форматированный адрес
	Entrance         null.String // Подъезд (опционально)
	Floor            null.String // Этаж (опционально)
	Apartment        null.String // Квартира (опционально)
	Comment          null.String // Комментарий к доставке
}

func (a *Address) ConvertToGRPCAddress() *gen.Address {
	return &gen.Address{
		Id:               a.ID.String(),
		Latitude:         a.Latitude,
		Longitude:        a.Longitude,
		FormattedAddress: a.FormattedAddress,
		Entrance:         nullableToString(a.Entrance),
		Floor:            nullableToString(a.Floor),
		Apartment:        nullableToString(a.Apartment),
		Comment:          nullableToString(a.Comment),
	}
}

func nullableToString(ns null.String) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
