package main

import (
	"2025_CakeLand_API/internal/pkg/config"
	handler2 "2025_CakeLand_API/internal/pkg/notification/delivery/grpc"
	"2025_CakeLand_API/internal/pkg/notification/delivery/grpc/generated"
	repo2 "2025_CakeLand_API/internal/pkg/notification/repo"
	"2025_CakeLand_API/internal/pkg/utils"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"2025_CakeLand_API/internal/pkg/utils/logger"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
func run() error {
	// Создаём Configuration
	conf, err := config.NewConfig()
	if err != nil {
		return err
	}
	// Создаём Logger
	l := logger.NewLogger(conf.Env)
	// Подключаем базу данных
	db, err := utils.ConnectPostgres(&conf.DB)
	if err != nil {
		return err
	}

	// Создаём grpc сервис
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPC.NotificationPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.LoggingUnaryInterceptor(l)),
	)

	tokenator := jwt.NewTokenator()
	mdProvider := md.NewMetadataProvider()
	repo := repo2.NewNotificationRepository(db)
	handler := handler2.NewNotificationHandler(l, mdProvider, repo, tokenator)

	generated.RegisterNotificationServiceServer(grpcServer, handler)
	l.Info("Starting notification gRPC server",
		slog.String("port", fmt.Sprintf(":%d", conf.GRPC.NotificationPort)),
	)
	return grpcServer.Serve(listener)
}
