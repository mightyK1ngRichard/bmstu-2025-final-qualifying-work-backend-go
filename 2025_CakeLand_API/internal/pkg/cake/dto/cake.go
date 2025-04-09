package dto

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type PreviewCake struct {
	ID              uuid.UUID   `db:"id"`
	Name            string      `db:"name"`
	PreviewImageURL string      `db:"image_url"`
	KgPrice         float64     `db:"kg_price"`
	Rating          uint        `db:"rating"`
	Description     null.String `db:"description"`
	Mass            float64     `db:"mass"`
	DiscountKgPrice null.Float  `db:"discount_kg_price"`
	DiscountEndTime null.Time   `db:"discount_end_time"`
	DateCreation    time.Time   `db:"date_creation"`
	IsOpenForSale   bool        `db:"is_open_for_sale"`
	OwnerID         uuid.UUID   `db:"owner_id"`
}

func (pc PreviewCake) ConvertToGrpcModel() *generated.PreviewCake {
	var description *wrappers.StringValue
	if pc.Description.Valid {
		description = &wrappers.StringValue{Value: pc.Description.String}
	}

	var discountKgPrice *wrappers.DoubleValue
	if pc.DiscountKgPrice.Valid {
		discountKgPrice = &wrappers.DoubleValue{Value: pc.DiscountKgPrice.Float64}
	}

	var discountEndTime *timestamp.Timestamp
	if pc.DiscountEndTime.Valid {
		discountEndTime = timestamppb.New(pc.DiscountEndTime.Time)
	}

	return &generated.PreviewCake{
		Id:              pc.ID.String(),
		Name:            pc.Name,
		PreviewImageUrl: pc.PreviewImageURL,
		KgPrice:         pc.KgPrice,
		Rating:          uint32(pc.Rating),
		Description:     description,
		Mass:            pc.Mass,
		DiscountKgPrice: discountKgPrice,
		DiscountEndTime: discountEndTime,
		DateCreation:    timestamppb.New(pc.DateCreation),
		IsOpenForSale:   pc.IsOpenForSale,
		OwnerId:         pc.OwnerID.String(),
	}
}
