package main

import (
	"2025_CakeLand_API/internal/pkg/config"
	"2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	handler "2025_CakeLand_API/internal/pkg/reviews/delivery/grpc"
	gen "2025_CakeLand_API/internal/pkg/reviews/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/reviews/repo"
	"2025_CakeLand_API/internal/pkg/reviews/usecase"
	"2025_CakeLand_API/internal/pkg/utils"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"2025_CakeLand_API/internal/pkg/utils/logger"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

// go run cmd/reviews/main.go --config=./config/config.yaml
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

	// Клиент сервиса пользователя
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", conf.GRPC.ProfilePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()
	userClient := generated.NewProfileServiceClient(conn)

	// Создаём grpc сервис
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPC.ReviewsPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.LoggingUnaryInterceptor(l)),
	)

	repository := repo.NewReviewsRepository(db)
	tokenator := jwt.NewTokenator()
	mdProvider := md.NewMetadataProvider()
	usecase := usecase.NewReviewsUsecase(userClient, repository)
	handler := handler.NewReviewsHandler(l, usecase, mdProvider, tokenator)
	gen.RegisterReviewServiceServer(grpcServer, handler)
	l.Info("Starting profile gRPC service", slog.String("port", fmt.Sprintf(":%d", conf.GRPC.ProfilePort)))
	return grpcServer.Serve(listener)
}
