package jwt

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

const (
	refreshTokenLifeSpan = 7 * 24 * time.Hour // Срок действия 7 дней
	accessTokenLifeSpan  = 15 * time.Minute   // Срок действия 15 минут
)

type Tokenator struct {
	accessSign  []byte
	refreshSign []byte
}

func NewTokenator() *Tokenator {
	return &Tokenator{
		accessSign:  []byte(os.Getenv("ACCESS_SIGN")),
		refreshSign: []byte(os.Getenv("REFRESH_SIGN")),
	}
}

// GenerateAccessToken генерирует access токен
func (t *Tokenator) GenerateAccessToken(userUID string) (*models.JWTTokenPayload, error) {
	return generateToken(userUID, accessTokenLifeSpan, t.accessSign)
}

// GenerateRefreshToken генерирует refresh токен
func (t *Tokenator) GenerateRefreshToken(userUID string) (*models.JWTTokenPayload, error) {
	return generateToken(userUID, refreshTokenLifeSpan, t.refreshSign)
}

// IsTokenExpired проверяет, истёк ли срок действия токена
func (t *Tokenator) IsTokenExpired(tokenString string, isRefresh bool) (bool, error) {
	var sign []byte
	if isRefresh {
		sign = t.refreshSign
	} else {
		sign = t.accessSign
	}

	// Извлечение claims и валидация токена
	claims, err := getTokenClaims(tokenString, sign)
	if err != nil {
		return false, err
	}

	if exp, ok := claims[string(domains.KeyExpClaim)].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		// Проверяем, истёк ли токен
		if time.Now().After(expirationTime) {
			return true, nil
		}
		return false, nil
	}

	return false, errs.ErrClaimIsMissing
}

// GetUserIDFromToken возвращает user_id если токен ещё не протух
func (t *Tokenator) GetUserIDFromToken(tokenString string, isRefresh bool) (string, error) {
	var sign []byte
	if isRefresh {
		sign = t.refreshSign
	} else {
		sign = t.accessSign
	}

	// Извлечение claims и валидация токена
	claims, err := getTokenClaims(tokenString, sign)
	if err != nil {
		return "", err
	}

	// Получает exp
	exp, ok := claims[domains.KeyExpClaim.String()].(float64)
	if !ok {
		return "", fmt.Errorf("%v: %s is missing in token", errs.ErrClaimIsMissing, domains.KeyExpClaim.String())
	}

	// Проверяем истёк ли токен
	expirationTime := time.Unix(int64(exp), 0)
	if time.Now().After(expirationTime) {
		return "", errs.ErrTokenIsExpired
	}

	// Достаём userID если токен не протух
	userID, ok := claims[domains.KeyUserIDClaim.String()].(string)
	if !ok {
		return "", fmt.Errorf("%v: %s is missing in token", errs.ErrClaimIsMissing, domains.KeyUserIDClaim.String())
	}

	return userID, nil
}

func generateToken(userUID string, duration time.Duration, sign []byte) (*models.JWTTokenPayload, error) {
	tokenExpiryTime := time.Now().Add(duration)
	claims := jwt.MapClaims{
		domains.KeyUserIDClaim.String(): userUID,
		domains.KeyExpClaim.String():    tokenExpiryTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sign)
	if err != nil {
		return nil, err
	}

	return &models.JWTTokenPayload{
		UserUID:   userUID,
		Token:     tokenString,
		ExpiresIn: tokenExpiryTime,
	}, nil
}

func getTokenClaims(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrUnexpectedSignInMethod
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrParsingToken, err)
	}

	// Извлечение claims и валидация токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errs.ErrInvalidTokenOrClaims
	}

	return claims, nil
}
