package models

import (
	"2025_CakeLand_API/internal/models/errs"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

var (
	ErrUserNotFound         = status.Error(codes.NotFound, "пользователь не найден")
	ErrInternal             = status.Error(codes.Internal, "ошибка на сервере")
	ErrPreviewImageNotFound = status.Error(codes.Internal, "ошибка создания preview")
	ErrNoToken              = status.Error(codes.Unauthenticated, "токен не найден")
	ErrNoMetadata           = status.Error(codes.InvalidArgument, "отсутствует metadata")
	ErrMissingFingerprint   = status.Error(codes.InvalidArgument, "отсутствует fingerprint")
	ErrInvalidPassword      = status.Error(codes.Unauthenticated, "неверный логин или пароль")
	ErrInvalidRefreshToken  = status.Error(codes.Unauthenticated, "неверный refresh токен")
	ErrUserAlreadyExists    = status.Error(codes.AlreadyExists, "пользователь с таким email уже существует")
	ErrTokenIsExpired       = status.Error(codes.Unauthenticated, "токен устарел")
	ErrExpMissingInToken    = status.Error(codes.Unauthenticated, "expiration time (exp) not found in token")
	ErrUserIDMissingInToken = status.Error(codes.Unauthenticated, "userID not found in token")
)

/* ################ DataBaseError ################ */

// DataBaseError ошибки для работы с базой данных.
type DataBaseError struct {
	Method string // Описание места ошибки
	Err    error  // Оригинальная ошибка
}

// NewDataBaseError Новый конструктор для DataBaseError
func NewDataBaseError(method string, err error) *DataBaseError {
	return &DataBaseError{
		Method: method,
		Err:    err,
	}
}

// Error Реализация интерфейса error
func (e *DataBaseError) Error() string {
	return fmt.Sprintf("Database error occurred in method %s with: %v", e.Method, e.Err)
}

// Unwrap Можем добавить дополнительный метод для извлечения оригинальной ошибки
func (e *DataBaseError) Unwrap() error {
	return e.Err
}

/* ################ ImageStorageError ################ */

// ImageStorageError - ошибка, возникающая при работе с хранилищем изображений.
type ImageStorageError struct {
	Message string
	Err     error
}

// Error implements the error interface for ImageStorageError.
func (e *ImageStorageError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Image storage error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("Image storage error: %s", e.Message)
}

// NewImageStorageError создает новую ошибку ImageStorageError.
func NewImageStorageError(message string, err error) *ImageStorageError {
	return &ImageStorageError{
		Message: message,
		Err:     err,
	}
}

// HandleError обрабатывает ошибку и возвращает gRPC-статус с подходящим кодом и логированием.
func HandleError(ctx context.Context, log *slog.Logger, err error, description string) error {
	if err == nil {
		return nil
	}

	switch {
	case isDatabaseError(err):
		logGRPCError(ctx, log, "database_error", err, description)
		return status.Error(codes.Internal, err.Error())

	case isImageStorageError(err):
		logGRPCError(ctx, log, "image_storage_error", err, description)
		return status.Error(codes.Internal, err.Error())

	case errors.Is(err, errs.ErrNotFound):
		logGRPCError(ctx, log, "not_found", err, description)
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, errs.ErrInvalidUUIDFormat):
		logGRPCError(ctx, log, "invalid_uuid_format", err, description)
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, errs.ErrUnexpectedSignInMethod):
		logGRPCError(ctx, log, "unexpected_signing_method", err, description)
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errs.ErrInvalidTokenOrClaims):
		logGRPCError(ctx, log, "invalid_token_or_claims", err, description)
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errs.ErrParsingToken):
		logGRPCError(ctx, log, "token_parsing", err, description)
		return status.Error(codes.Unauthenticated, err.Error())
	}

	// Если это уже gRPC-ошибка — вернём как есть
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Иначе — неизвестная ошибка
	logGRPCError(ctx, log, "unknown_error", err, description)
	return status.Errorf(codes.Unknown, "Неизвестная ошибка: %v", err.Error())
}

func isDatabaseError(err error) bool {
	var dbErr *DataBaseError
	return errors.As(err, &dbErr)
}

func isImageStorageError(err error) bool {
	var imgErr *ImageStorageError
	return errors.As(err, &imgErr)
}

// logGRPCError логирует gRPC-ошибку с единообразным форматом.
func logGRPCError(ctx context.Context, log *slog.Logger, kind string, err error, description string) {
	log.Log(
		ctx,
		slog.LevelWarn,
		"grpc error",
		slog.String("type", kind),
		slog.String("description", description),
		slog.String("error", err.Error()),
	)
}
