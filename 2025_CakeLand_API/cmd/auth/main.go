package main

import (
	auth "2025_CakeLand_API/internal/pkg/auth/delivery/grpc"
	"2025_CakeLand_API/internal/pkg/auth/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/auth/repo"
	"2025_CakeLand_API/internal/pkg/auth/usecase"
	"2025_CakeLand_API/internal/pkg/config"
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

// go run cmd/auth/main.go --config=./config/config.yaml
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
	log := logger.NewLogger(conf.Env)
	// Подключаем базу данных
	db, err := utils.ConnectPostgres(&conf.DB)
	if err != nil {
		return err
	}

	// Создаём grpc сервис
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPC.AuthPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	rep := repo.NewAuthRepository(db)
	validator := utils.NewValidator()
	tokenator := jwt.NewTokenator()
	mdProvider := md.NewMetadataProvider()
	authUsecase := usecase.NewAuthUsecase(log, tokenator, rep)
	grpcAuthHandler := auth.NewGrpcAuthHandler(validator, authUsecase, mdProvider)

	generated.RegisterAuthServer(grpcServer, grpcAuthHandler)
	log.Info("Starting gRPC server", slog.String("port", fmt.Sprintf(":%d", conf.GRPC.AuthPort)))
	return grpcServer.Serve(listener)
}
