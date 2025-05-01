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
	ReviewsCount    uint
	StarsSum        uint
	Description     null.String
	Mass            float64
	DiscountKgPrice null.Float
	DiscountEndTime null.Time
	DateCreation    time.Time
	IsOpenForSale   bool
	Owner           Owner
	ColorsHex       []string
}

type PreviewCakeDB struct {
	ID              uuid.UUID
	Name            string
	PreviewImageURL string
	KgPrice         float64
	ReviewsCount    uint
	StarsSum        uint
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

	// рейтинг = сумма / кол-во
	var rating uint32 = 0
	if pc.ReviewsCount != 0 {
		rating = uint32(pc.StarsSum / pc.ReviewsCount)
	}

	return &generated.PreviewCake{
		Id:              pc.ID.String(),
		Name:            pc.Name,
		PreviewImageUrl: pc.PreviewImageURL,
		KgPrice:         pc.KgPrice,
		Rating:          rating,
		Description:     description,
		Mass:            pc.Mass,
		DiscountKgPrice: discountKgPrice,
		DiscountEndTime: discountEndTime,
		DateCreation:    timestamppb.New(pc.DateCreation),
		IsOpenForSale:   pc.IsOpenForSale,
		Owner:           pc.Owner.ConvertToGrpcUser(),
		ColorsHex:       pc.ColorsHex,
	}
}

func (pc *PreviewCakeDB) ConvertToPreviewCake(owner Owner) PreviewCake {
	return PreviewCake{
		ID:              pc.ID,
		Name:            pc.Name,
		PreviewImageURL: pc.PreviewImageURL,
		KgPrice:         pc.KgPrice,
		StarsSum:        pc.StarsSum,
		ReviewsCount:    pc.ReviewsCount,
		Description:     pc.Description,
		Mass:            pc.Mass,
		DiscountKgPrice: pc.DiscountKgPrice,
		DiscountEndTime: pc.DiscountEndTime,
		DateCreation:    pc.DateCreation,
		IsOpenForSale:   pc.IsOpenForSale,
		Owner:           owner,
	}
}
