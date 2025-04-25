package logger

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		duration := time.Since(start)

		p, _ := peer.FromContext(ctx)
		st, _ := status.FromError(err)

		logger.Info("gRPC request",
			slog.String("method", info.FullMethod),
			slog.String("peer", p.Addr.String()),
			slog.Duration("duration", duration),
			slog.Any("error", err),
			slog.String("code", st.Code().String()),
		)
		return resp, err
	}
}
