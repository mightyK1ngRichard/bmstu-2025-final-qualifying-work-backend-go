package main

import (
	cake "2025_CakeLand_API/internal/pkg/cake/delivery/grpc"
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/cake/repo"
	"2025_CakeLand_API/internal/pkg/cake/usecase"
	"2025_CakeLand_API/internal/pkg/config"
	minioStorage "2025_CakeLand_API/internal/pkg/s3storage"
	"2025_CakeLand_API/internal/pkg/utils"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"2025_CakeLand_API/internal/pkg/utils/logger"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"fmt"
	"log/slog"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

// go run cmd/cake/main.go --config=./config/config.yaml
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

	// Создаём S3 хранилище
	minioClient, err := minioStorage.NewMinioClient(&conf.MinIO)
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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPC.Port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024), // 10MB для входящих сообщений
		grpc.MaxSendMsgSize(10*1024*1024), // 10MB для исходящих сообщений
	)
	repository := repo.NewCakeRepository(db)
	tokenator := jwt.NewTokenator()
	useCase := usecase.NewCakeUsecase(l, tokenator, repository, minioClient, conf.MinIO.Bucket)
	mdProvider := md.NewMetadataProvider()
	handler := cake.NewCakeHandler(l, useCase, mdProvider)
	generated.RegisterCakeServiceServer(grpcServer, handler)
	l.Info("Starting gRPC cake service", slog.String("port", fmt.Sprintf(":%d", conf.GRPC.Port)))
	return grpcServer.Serve(listener)
}
