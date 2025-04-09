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

// HandleError обрабатывает ошибку и возвращает соответствующую gRPC ошибку с нужным кодом и сообщением.
func HandleError(ctx context.Context, log *slog.Logger, err error, description string) error {
	if err != nil {
		// Обработка ошибки базы данных
		var dbErr *DataBaseError
		if errors.As(err, &dbErr) {
			// Логируем ошибку с уровнем предупреждения
			log.Log(ctx, slog.LevelWarn, "Database error", slog.String("description", description), slog.String("error", err.Error()))
			return status.Error(codes.Internal, dbErr.Error())
		}

		// Обработка ошибки хранилища изображений
		var imgErr *ImageStorageError
		if errors.As(err, &imgErr) {
			// Логируем ошибку с уровнем предупреждения
			log.Log(ctx, slog.LevelWarn, "Image store error", slog.String("description", description), slog.String("error", err.Error()))
			return status.Error(codes.Internal, imgErr.Error())
		}

		switch {
		case errors.Is(err, errs.ErrNotFound):
			log.Log(ctx, slog.LevelWarn, "Not found error", slog.String("description", description), slog.String("error", err.Error()))
			return status.Error(codes.NotFound, err.Error())
		case errors.Is(err, errs.ErrInvalidUUIDFormat):
			log.Log(ctx, slog.LevelWarn, "Invalid UUID format", slog.String("description", description), slog.String("error", err.Error()))
			return status.Error(codes.InvalidArgument, err.Error())
		}

		// Проверка на стандартные ошибки
		if st, ok := status.FromError(err); ok {
			return st.Err()
		}

		// Логируем неизвестную ошибку
		log.Log(ctx, slog.LevelWarn, "Unknown error", slog.String("description", description), slog.String("error", err.Error()))
		return status.Errorf(codes.Unknown, "Неизвестная ошибка: %v", err.Error())
	}

	return nil
}
