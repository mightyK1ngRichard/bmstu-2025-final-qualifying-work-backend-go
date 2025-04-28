package main

import (
	"2025_CakeLand_API/internal/pkg/config"
	"2025_CakeLand_API/internal/pkg/order/delivery/grpc"
	"2025_CakeLand_API/internal/pkg/order/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/order/repo"
	"2025_CakeLand_API/internal/pkg/order/usecase"
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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPC.OrderPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.LoggingUnaryInterceptor(l)),
		grpc.MaxRecvMsgSize(20*1024*1024), // 200MB для входящих сообщений
		grpc.MaxSendMsgSize(20*1024*1024), // 200MB для исходящих сообщений
	)
	repository := repo.NewOrderRepo(db)
	tokenator := jwt.NewTokenator()
	uc := usecase.NewOrderUsecase(tokenator, repository)
	mdProvider := md.NewMetadataProvider()
	h := handler.NewOrderHandler(l, uc, mdProvider)
	generated.RegisterOrderServiceServer(grpcServer, h)
	l.Info("Starting order gRPC service", slog.String("port", fmt.Sprintf(":%d", conf.GRPC.OrderPort)))
	return grpcServer.Serve(listener)
}
