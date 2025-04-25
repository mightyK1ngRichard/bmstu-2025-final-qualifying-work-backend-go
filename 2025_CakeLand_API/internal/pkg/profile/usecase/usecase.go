package usecase

import (
	"2025_CakeLand_API/internal/models"
	dto2 "2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/minio"
	"2025_CakeLand_API/internal/pkg/profile"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"context"
	"github.com/google/uuid"
	"sync"
)

type ProfileUseсase struct {
	tokenator     *jwt.Tokenator
	repo          profile.IProfileRepository
	imageProvider *minio.MinioProvider
}

func NewProfileUsecase(
	tokenator *jwt.Tokenator,
	repo profile.IProfileRepository,
	imageProvider *minio.MinioProvider,
) *ProfileUseсase {
	return &ProfileUseсase{
		tokenator:     tokenator,
		repo:          repo,
		imageProvider: imageProvider,
	}
}

func (u *ProfileUseсase) UserInfo(ctx context.Context, accessToken string) (*dto.UserInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Достаём UserID
	userIDStr, err := u.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	//var userInfo dto.UserInfo
	var (
		dbCakes []dto2.PreviewCakeDB
		user    *dto.Profile
	)
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	wg.Add(2)

	// Получаем данные пользователя
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		profileDB, err := u.repo.UserInfo(ctx, userID)
		if err != nil {
			trySendError(err, errChan, cancel)
			return
		}

		if ctx.Err() != nil {
			return
		}

		user = profileDB
	}()

	// Получаем тортики пользователя
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		dbPreviewCakes, err := u.repo.CakesByUserID(ctx, userID)
		if err != nil {
			trySendError(err, errChan, cancel)
			return
		}

		if ctx.Err() != nil {
			return
		}

		dbCakes = dbPreviewCakes
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, err
	}

	cakes := make([]dto2.PreviewCake, len(dbCakes))
	owner := user.ConvertToOwner()
	for i, cake := range dbCakes {
		cakes[i] = cake.ConvertToPreviewCake(owner)
	}
	userInfo := dto.UserInfo{
		User:  *user,
		Cakes: cakes,
	}
	return &userInfo, nil
}

func (u *ProfileUseсase) UserInfoByID(ctx context.Context, userID uuid.UUID) (*models.UserInfo, error) {
	profileInfo, err := u.repo.UserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo := models.UserInfo{
		ID:             profileInfo.ID.String(),
		FIO:            profileInfo.FIO,
		Address:        profileInfo.Address,
		Nickname:       profileInfo.Nickname,
		ImageURL:       profileInfo.ImageURL,
		HeaderImageURL: profileInfo.HeaderImageURL,
		Mail:           profileInfo.Mail,
		Phone:          profileInfo.Phone,
	}
	return &userInfo, nil
}

func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}
