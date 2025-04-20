package reviews

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/reviews/entities"
	"context"
	"github.com/google/uuid"
)

type IReviewsUsecase interface {
	CreateFeedback(context.Context, entities.CreateFeedbackReq) (*models.Feedback, error)
	ProductFeedbacks(context.Context, uuid.UUID) ([]models.Feedback, error)
}

type IReviewsRepository interface {
	AddFeedback(context.Context, *models.FeedbackDB) error
	ProductFeedbacks(context.Context, uuid.UUID) ([]models.FeedbackDB, error)
}
