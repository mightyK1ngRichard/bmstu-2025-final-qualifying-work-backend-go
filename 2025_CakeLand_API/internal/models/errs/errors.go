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
)

func ConvertToGrpcError(ctx context.Context, log *slog.Logger, err error, description string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrPreviewImageNotFound):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.Internal, fmt.Sprintf("%s: %v", description, err))

	case errors.Is(err, ErrTokenIsExpired):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.Unauthenticated, fmt.Sprintf("%s: %v", description, err))

	case errors.Is(err, ErrClaimIsMissing):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.Internal, fmt.Sprintf("%s: %v", description, err))

	case errors.Is(err, ErrNotFound):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.NotFound, fmt.Sprintf("%s: %v", description, err))

	case errors.Is(err, ErrAlreadyExists):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.AlreadyExists, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrNoToken):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.Unauthenticated, "missing token")

	case errors.Is(err, ErrInvalidPassword):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.InvalidArgument, "invalid email or password")

	case errors.Is(err, ErrNoMetadata):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrInvalidUUIDFormat),
		errors.Is(err, ErrInvalidRefreshToken):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %s", err, description))

	case errors.Is(err, ErrUnexpectedSignInMethod),
		errors.Is(err, ErrInvalidTokenOrClaims),
		errors.Is(err, ErrParsingToken):
		logGRPCError(ctx, log, err, description)
		return status.Error(codes.Unauthenticated, fmt.Sprintf("%v: %s", err, description))
	}

	// Если это уже gRPC-ошибка — вернём как есть
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Неизвестная ошибка
	logGRPCError(ctx, log, err, description)
	return status.Errorf(codes.Unknown, "%s: %v", description, err)
}

func logGRPCError(ctx context.Context, log *slog.Logger, err error, description string) {
	log.Log(
		ctx,
		slog.LevelWarn,
		"grpc error",
		slog.String("description", description),
		slog.String("error", err.Error()),
	)
}
