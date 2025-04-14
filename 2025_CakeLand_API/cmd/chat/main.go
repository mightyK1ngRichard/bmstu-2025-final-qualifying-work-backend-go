package main

import (
	chat "2025_CakeLand_API/internal/pkg/chat/delivery/grpc"
	"2025_CakeLand_API/internal/pkg/chat/delivery/grpc/generated"
	chatRepo "2025_CakeLand_API/internal/pkg/chat/repo"
	"2025_CakeLand_API/internal/pkg/config"
	"2025_CakeLand_API/internal/pkg/utils"
	"2025_CakeLand_API/internal/pkg/utils/logger"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
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

	// Подключаем базу данных
	db, err := utils.ConnectPostgres(&conf.DB)
	if err != nil {
		return err
	}

	chatPort := fmt.Sprintf(":%d", conf.GRPC.ChatPort)
	lis, err := net.Listen("tcp", chatPort)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	// Создаём Logger
	l := logger.NewLogger(conf.Env)

	grpcServer := grpc.NewServer()
	repo := chatRepo.NewChatRepository(db)
	chatProvider := chat.NewChatProvider(l, repo)
	generated.RegisterChatServiceServer(grpcServer, chatProvider)

	l.Info("Starting chat gRPC service", slog.String("port", chatPort))
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
		return err
	}

	return nil
}
