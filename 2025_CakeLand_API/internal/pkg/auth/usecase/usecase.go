package usecase

import (
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/auth"
	"2025_CakeLand_API/internal/pkg/auth/dto"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AuthUseсase struct {
	tokenator *jwt.Tokenator
	repo      auth.IAuthRepository
}

func NewAuthUsecase(
	tokenator *jwt.Tokenator,
	repo auth.IAuthRepository,
) *AuthUseсase {
	return &AuthUseсase{
		tokenator: tokenator,
		repo:      repo,
	}
}

func (u *AuthUseсase) Login(ctx context.Context, in dto.LoginReq) (*dto.LoginRes, error) {
	// Получаем данные пользователя
	res, err := u.repo.GetUserByEmail(ctx, dto.GetUserByEmailReq{
		Email: in.Email,
	})
	if err != nil {
		return nil, err
	}

	// Проверяем пароль пользователя
	if !checkPassword(in.Password, res.PasswordHash) {
		return nil, errs.ErrInvalidPassword
	}

	// Создаём новый access токен
	accessToken, err := u.tokenator.GenerateAccessToken(res.ID.String())
	if err != nil {
		return nil, err
	}

	oldRefreshToken, exists := res.RefreshTokensMap[in.Fingerprint]
	// Если токен уже существует, проверим его срок годности. Если ещё валиден, тогда генерируем только access token
	if exists {
		isExpired, expErr := u.tokenator.IsTokenExpired(oldRefreshToken, true)
		if expErr != nil {
			// Если не вышло декодировать токен, создадим новый токен
		} else if !isExpired {
			// Если токен не устарел, создаём только access токен
			return &dto.LoginRes{
				AccessToken:  accessToken.Token,
				RefreshToken: oldRefreshToken,
				ExpiresIn:    accessToken.ExpiresIn,
			}, nil
		}
	}

	// Создаём новый refresh токен
	newRefreshToken, err := u.tokenator.GenerateRefreshToken(res.ID.String())
	if err != nil {
		return nil, err
	}

	// Сохраняем или обновляем токены в бд
	res.RefreshTokensMap[in.Fingerprint] = newRefreshToken.Token
	err = u.repo.UpdateUserRefreshTokens(ctx, dto.UpdateUserRefreshTokensReq{
		UserID:           res.ID.String(),
		RefreshTokensMap: res.RefreshTokensMap,
	})
	if err != nil {
		return nil, err
	}

	return &dto.LoginRes{
		AccessToken:  accessToken.Token,
		RefreshToken: newRefreshToken.Token,
		ExpiresIn:    accessToken.ExpiresIn,
	}, nil
}

func (u *AuthUseсase) Register(ctx context.Context, in dto.RegisterReq) (*dto.RegisterRes, error) {
	hashedPassword, err := generatePasswordHash(in.Password)
	if err != nil {
		return nil, err
	}

	// Создаём токены
	userID := uuid.New()
	accessToken, errAccess := u.tokenator.GenerateAccessToken(userID.String())
	refreshToken, errRefresh := u.tokenator.GenerateRefreshToken(userID.String())
	if errAccess != nil {
		return nil, errAccess
	} else if errRefresh != nil {
		return nil, errRefresh
	}

	// Создаём пользователя
	if err = u.repo.CreateUser(ctx, dto.CreateUserReq{
		UUID:         userID,
		Email:        strings.TrimSpace(in.Email),
		Nickname:     strings.TrimSpace(in.Nickname),
		PasswordHash: hashedPassword,
		RefreshTokensMap: map[string]string{
			in.Fingerprint: refreshToken.Token,
		},
	}); err != nil {
		return nil, err
	}

	return &dto.RegisterRes{
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresIn:    accessToken.ExpiresIn,
	}, nil
}

func (u *AuthUseсase) UpdateAccessToken(ctx context.Context, in dto.UpdateAccessTokenReq) (*dto.UpdateAccessTokenRes, error) {
	// Получаем userID пользователя из refresh токена
	userID, err := u.tokenator.GetUserIDFromToken(in.RefreshToken, true)
	if err != nil {
		return nil, err
	}

	// Получаем все refresh токены пользователя
	res, err := u.repo.GetUserRefreshTokens(ctx, dto.GetUserRefreshTokensReq{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	// Ищем refresh токен для заданного fingerprint
	oldRefreshToken, exists := res.RefreshTokensMap[in.Fingerprint]
	if !exists {
		return nil, errs.ErrNoToken
	}

	// Проверяем схожи ли токены
	if oldRefreshToken != in.RefreshToken {
		return nil, errs.ErrInvalidRefreshToken
	}

	// Генерируем новый access токен
	accessToken, err := u.tokenator.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateAccessTokenRes{
		AccessToken: accessToken.Token,
		ExpiresIn:   accessToken.ExpiresIn,
	}, nil
}

func (u *AuthUseсase) Logout(ctx context.Context, in dto.LogoutReq) (*dto.LogoutRes, error) {
	// Получение userID из refresh токена
	userID, err := u.tokenator.GetUserIDFromToken(in.RefreshToken, true)
	if err != nil {
		return nil, err
	}

	// Получаем токены пользователя
	res, err := u.repo.GetUserRefreshTokens(ctx, dto.GetUserRefreshTokensReq{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	// Проверяем, верный ли refresh токен
	dbRefreshToken := res.RefreshTokensMap[in.Fingerprint]
	if dbRefreshToken != in.RefreshToken {
		return nil, errs.ErrInvalidRefreshToken
	}

	delete(res.RefreshTokensMap, in.Fingerprint)
	err = u.repo.UpdateUserRefreshTokens(ctx, dto.UpdateUserRefreshTokensReq{
		UserID:           userID,
		RefreshTokensMap: res.RefreshTokensMap,
	})
	if err != nil {
		return nil, err
	}

	return &dto.LogoutRes{
		Message: "Logged out successfully",
	}, nil
}

func generatePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}

func checkPassword(inputPassword string, realPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(realPassword, []byte(inputPassword))
	return err == nil
}
