package usecase

import (
	"2025_CakeLand_API/internal/pkg/minio"
	"2025_CakeLand_API/internal/pkg/profile"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"context"
	"github.com/google/uuid"
	"log/slog"
	"sync"
)

type ProfileUseсase struct {
	log           *slog.Logger
	tokenator     *jwt.Tokenator
	repo          profile.IProfileRepository
	imageProvider *minio.MinioProvider
}

func NewProfileUsecase(
	log *slog.Logger,
	tokenator *jwt.Tokenator,
	repo profile.IProfileRepository,
	imageProvider *minio.MinioProvider,
) *ProfileUseсase {
	return &ProfileUseсase{
		log:           log,
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

	var userInfo dto.UserInfo
	mu := sync.Mutex{}
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

		mu.Lock()
		userInfo.User = *profileDB
		mu.Unlock()
	}()

	// Получаем тортики пользователя
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		previewCakes, err := u.repo.CakesByUserID(ctx, userID)
		if err != nil {
			trySendError(err, errChan, cancel)
			return
		}

		if ctx.Err() != nil {
			return
		}

		mu.Lock()
		userInfo.Cakes = previewCakes
		mu.Unlock()
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, err
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
