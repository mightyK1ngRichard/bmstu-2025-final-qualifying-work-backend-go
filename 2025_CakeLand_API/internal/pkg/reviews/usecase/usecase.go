package usecase

import (
	"2025_CakeLand_API/internal/models"
	profileGen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/reviews"
	"2025_CakeLand_API/internal/pkg/reviews/entities"
	"context"
	"github.com/google/uuid"
	"sync"
)

type ReviewsUseсase struct {
	repo          reviews.IReviewsRepository
	profileClient profileGen.ProfileServiceClient
}

func NewReviewsUsecase(
	profileClient profileGen.ProfileServiceClient,
	repo reviews.IReviewsRepository,
) *ReviewsUseсase {
	return &ReviewsUseсase{
		repo:          repo,
		profileClient: profileClient,
	}
}

func (u *ReviewsUseсase) CreateFeedback(ctx context.Context, req entities.CreateFeedbackReq) (*models.Feedback, error) {
	// Создаём отзыв в бд
	dbFeedback := models.FeedbackDB{
		ID:       uuid.New(),
		Text:     req.Text,
		AuthorID: req.AuthorID,
		CakeID:   req.CakeID,
	}
	err := u.repo.AddFeedback(ctx, &dbFeedback)
	if err != nil {
		return nil, err
	}

	// Получаем данные пользователя
	userInfo := models.UserInfo{}
	res, err := u.profileClient.GetUserInfoByID(ctx, &profileGen.GetUserInfoByIDReq{
		UserID: dbFeedback.AuthorID.String(),
	})
	if err == nil {
		if user := models.NewUserInfo(res.User); user != nil {
			userInfo = *user
		}
	}

	feedback := dbFeedback.ConvertToFeedback(userInfo)
	return &feedback, err
}

func (u *ReviewsUseсase) ProductFeedbacks(ctx context.Context, productID uuid.UUID) ([]models.Feedback, error) {
	dbFeedbacks, err := u.repo.ProductFeedbacks(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Получаем данные по авторам
	feedbacks := make([]models.Feedback, len(dbFeedbacks))
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}
	for i, feedback := range dbFeedbacks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			res, errDB := u.profileClient.GetUserInfoByID(ctx, &profileGen.GetUserInfoByIDReq{
				UserID: feedback.AuthorID.String(),
			})
			if errDB != nil {
				return
			}

			user := models.NewUserInfo(res.User)
			if user == nil {
				return
			}

			mu.Lock()
			feedbacks[i] = feedback.ConvertToFeedback(*user)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return feedbacks, nil
}
