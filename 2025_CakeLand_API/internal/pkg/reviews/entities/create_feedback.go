package entities

import (
	"2025_CakeLand_API/internal/models/errs"
	gen "2025_CakeLand_API/internal/pkg/reviews/delivery/grpc/generated"
	"fmt"
	"github.com/google/uuid"
)

type CreateFeedbackReq struct {
	Text     string
	Rating   int
	CakeID   uuid.UUID
	AuthorID uuid.UUID
}

func NewCreateFeedbackReq(req *gen.AddFeedbackRequest, authorID string) (*CreateFeedbackReq, error) {
	cakeID, err := uuid.Parse(req.GetCakeID())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err)
	}
	authorUID, err := uuid.Parse(authorID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUUIDFormat, err)
	}

	return &CreateFeedbackReq{
		Text:     req.GetText(),
		Rating:   int(req.GetRating()),
		CakeID:   cakeID,
		AuthorID: authorUID,
	}, nil
}
