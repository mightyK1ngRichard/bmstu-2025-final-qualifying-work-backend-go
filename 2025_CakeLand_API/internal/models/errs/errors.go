package errs

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrInvalidUUIDFormat      = errors.New("invalid uuid format")
	ErrUnexpectedSignInMethod = errors.New("unexpected signing method")
	ErrInvalidTokenOrClaims   = errors.New("invalid token or claims")
	ErrParsingToken           = errors.New("error parsing token")
	ErrAlreadyExists          = errors.New("already exists")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrNoToken                = errors.New("no token")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrNoMetadata             = errors.New("no metadata")
	ErrTokenIsExpired         = errors.New("token is expired")
	ErrClaimIsMissing         = errors.New("claim is missing")
	ErrPreviewImageNotFound   = errors.New("preview image not found")
	ErrDB                     = errors.New("database error")
	ErrNoMessage              = errors.New("no message")
	ErrInvalidInput           = errors.New("invalid input")
)

func ConvertToGrpcError(ctx context.Context, log *slog.Logger, err error, description string) error {
	if err == nil {
		return nil
	}

	logGRPCError(ctx, log, err, description)

	switch {
	case errors.Is(err, ErrDB):
		return status.Error(codes.Internal, "internal server error")

	case errors.Is(err, ErrNoMessage):
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrPreviewImageNotFound):
		return status.Error(codes.Internal, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrTokenIsExpired):
		return status.Error(codes.Unauthenticated, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrClaimIsMissing):
		return status.Error(codes.Internal, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrNotFound):
		return status.Error(codes.NotFound, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrNoToken):
		return status.Error(codes.Unauthenticated, "missing token")

	case errors.Is(err, ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, "invalid email or password")

	case errors.Is(err, ErrNoMetadata):
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrInvalidUUIDFormat),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrInvalidRefreshToken):
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrUnexpectedSignInMethod),
		errors.Is(err, ErrInvalidTokenOrClaims),
		errors.Is(err, ErrParsingToken):
		return status.Error(codes.Unauthenticated, fmt.Sprintf("%v: %s", err, description))
	}

	// Если это уже gRPC-ошибка — вернём как есть
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Неизвестная ошибка
	return status.Errorf(codes.Unknown, "unknown error")
}

func logGRPCError(ctx context.Context, log *slog.Logger, err error, description string) {
	log.Log(
		ctx,
		slog.LevelWarn,
		"error",
		slog.String("description", description),
		slog.String("error", err.Error()),
	)
}

func WrapDBError(method string, err error) error {
	return fmt.Errorf("%w: %s: %w", ErrDB, method, err)
}
