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
	queryProductFeedbacks   = `SELECT id, text, date_creation, rating, cake_id, author_id FROM feedback WHERE cake_id = $1`
	queryAddProductFeedback = `INSERT INTO feedback (id, text, rating, cake_id, author_id) VALUES ($1, $2, $3, $4, $5)`
)

type ReviewsRepository struct {
	db *sql.DB
}

func NewReviewsRepository(db *sql.DB) *ReviewsRepository {
	return &ReviewsRepository{
		db: db,
	}
}

func (r *ReviewsRepository) AddFeedback(ctx context.Context, feedback *models.FeedbackDB) error {
	const methodName = "[Repo.AddFeedback]"

	if _, err := r.db.ExecContext(ctx, queryAddProductFeedback,
		feedback.ID, feedback.Text, feedback.Rating, feedback.CakeID, feedback.AuthorID,
	); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *ReviewsRepository) ProductFeedbacks(ctx context.Context, id uuid.UUID) ([]models.FeedbackDB, error) {
	const methodName = "[Repo.ProductFeedbacks]"

	rows, err := r.db.QueryContext(ctx, queryProductFeedbacks, id)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var feedbacks []models.FeedbackDB
	for rows.Next() {
		var feedback models.FeedbackDB
		if err = rows.Scan(
			&feedback.ID,
			&feedback.Text,
			&feedback.DateCreation,
			&feedback.Rating,
			&feedback.CakeID,
			&feedback.AuthorID,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.ErrNotFound
			}
			return nil, errs.WrapDBError(methodName, err)
		}

		feedbacks = append(feedbacks, feedback)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return feedbacks, nil
}
