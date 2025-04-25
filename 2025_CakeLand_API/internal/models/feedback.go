package models

import (
	gen "2025_CakeLand_API/internal/pkg/reviews/delivery/grpc/generated"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Feedback struct {
	ID           uuid.UUID
	Text         string
	DateCreation time.Time
	Rating       int
	CakeID       uuid.UUID
	Author       UserInfo
}

type FeedbackDB struct {
	ID           uuid.UUID
	Text         string
	DateCreation time.Time
	Rating       int
	CakeID       uuid.UUID
	AuthorID     uuid.UUID
}

func (f *Feedback) ConvertToGRPC() *gen.Feedback {
	author := f.Author.ConvertToGRPCProfile()

	return &gen.Feedback{
		Id:           f.ID.String(),
		Text:         f.Text,
		DateCreation: timestamppb.New(f.DateCreation),
		Rating:       int32(f.Rating),
		CakeId:       f.CakeID.String(),
		Author:       author,
	}
}

func (f *FeedbackDB) ConvertToFeedback(author UserInfo) Feedback {
	return Feedback{
		ID:           f.ID,
		Text:         f.Text,
		DateCreation: f.DateCreation,
		Rating:       f.Rating,
		CakeID:       f.CakeID,
		Author:       author,
	}
}
