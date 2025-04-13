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
	ID              uuid.UUID
	Name            string
	PreviewImageURL string
	KgPrice         float64
	Rating          uint
	Description     null.String
	Mass            float64
	DiscountKgPrice null.Float
	DiscountEndTime null.Time
	DateCreation    time.Time
	IsOpenForSale   bool
	Owner           Owner
}

type PreviewCakeDB struct {
	ID              uuid.UUID
	Name            string
	PreviewImageURL string
	KgPrice         float64
	Rating          uint
	Description     null.String
	Mass            float64
	DiscountKgPrice null.Float
	DiscountEndTime null.Time
	DateCreation    time.Time
	IsOpenForSale   bool
	OwnerID         uuid.UUID
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
		Owner:           pc.Owner.ConvertToGrpcUser(),
	}
}

func (pc *PreviewCakeDB) ConvertToPreviewCake(owner Owner) PreviewCake {
	return PreviewCake{
		ID:              pc.ID,
		Name:            pc.Name,
		PreviewImageURL: pc.PreviewImageURL,
		KgPrice:         pc.KgPrice,
		Rating:          pc.Rating,
		Description:     pc.Description,
		Mass:            pc.Mass,
		DiscountKgPrice: pc.DiscountKgPrice,
		DiscountEndTime: pc.DiscountEndTime,
		DateCreation:    pc.DateCreation,
		IsOpenForSale:   pc.IsOpenForSale,
		Owner:           owner,
	}
}
