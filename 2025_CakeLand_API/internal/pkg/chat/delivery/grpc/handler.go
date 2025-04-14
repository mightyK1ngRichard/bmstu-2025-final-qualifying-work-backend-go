package grpc

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	gen "2025_CakeLand_API/internal/pkg/chat/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/chat/repo"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log/slog"
	"sync"
	"time"
)

type ChatProvider struct {
	gen.UnimplementedChatServiceServer
	clients    map[string]gen.ChatService_ChatServer
	mdProvider *md.MetadataProvider
	tokenator  *jwt.Tokenator
	log        *slog.Logger
	mu         sync.Mutex
	repo       repo.IChatRepository
}

func NewChatProvider(
	log *slog.Logger,
	mdProvider *md.MetadataProvider,
	tokenator *jwt.Tokenator,
	repo repo.IChatRepository,
) *ChatProvider {
	return &ChatProvider{
		clients:    make(map[string]gen.ChatService_ChatServer),
		mdProvider: mdProvider,
		tokenator:  tokenator,
		log:        log,
		repo:       repo,
	}
}

func (p *ChatProvider) Chat(stream gen.ChatService_ChatServer) error {
	ctx := stream.Context()

	// Получаем токен из метаданных
	accessToken, err := p.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return errs.ConvertToGrpcError(ctx, p.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	// Получаем UserID из токена
	ownerID, err := p.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return errs.ConvertToGrpcError(ctx, p.log, err, "failed to fetch user id from token")
	}

	var fromID string
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			p.removeClient(fromID)
			return nil
		}
		if err != nil {
			return errs.ConvertToGrpcError(stream.Context(), p.log, err, "error receiving message from server")
		}

		// Если это первое сообщение, запоминаем клиента
		if fromID == "" {
			fromID = ownerID
			p.addClient(fromID, stream)
		}

		// Если нет адресата, ничего не делаем
		if msg.ReceiverID == "" {
			continue
		}

		// Если время не указано, устанавливаем его
		var creationTime time.Time
		if msg.DateCreation == nil {
			creationTime = time.Now()
		} else {
			creationTime = msg.DateCreation.AsTime()
		}
		msg.DateCreation = timestamppb.New(creationTime)

		// Сохраняем в бд
		go func() {
			message := models.Message{
				ID:           uuid.NewString(),
				Text:         msg.Text,
				OwnerID:      ownerID,
				ReceiverID:   msg.ReceiverID,
				DateCreation: creationTime,
			}

			if err = p.repo.AddMessage(stream.Context(), message); err != nil {
				p.log.Warn("Error adding message to repo", "error", err)
			}
		}()

		p.mu.Lock()
		if client, ok := p.clients[msg.ReceiverID]; ok {
			client.Send(msg)
		} else {
			p.log.Warn(fromID, "client not found", msg.ReceiverID)
		}
		p.mu.Unlock()
	}
}

func (p *ChatProvider) UserChats(ctx context.Context, _ *emptypb.Empty) (*gen.UserChatsResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Получаем токен из метаданных
	accessToken, err := p.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, p.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	// Получаем UserID из токена
	userID, err := p.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, p.log, err, "failed to fetch user id from token")
	}

	// Получаем всех пользователей
	interlocutors, err := p.repo.UserInterlocutors(ctx, userID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, p.log, err, "failed to fetch interlocutors")
	}

	uniqueInterlocutors := uniqueStrings(interlocutors)

	// Получаем данные по пользователям
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	interlocutorsInfo := make([]*generated.User, len(interlocutors))
	for index, interlocutorID := range uniqueInterlocutors {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			interlocutor, userErr := p.repo.UserByID(ctx, interlocutorID)
			if userErr != nil {
				trySendError(userErr, errChan, cancel)
				return
			}

			mu.Lock()
			interlocutorsInfo[index] = interlocutor.ConvertToUserGRPC()
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, errs.ConvertToGrpcError(ctx, p.log, err, "failed to fetch interlocutors")
	}

	return &gen.UserChatsResponse{
		Users: interlocutorsInfo,
	}, nil
}

func (p *ChatProvider) ChatHistory(ctx context.Context, in *gen.ChatHistoryRequest) (*gen.ChatHistoryResponse, error) {
	// Получаем токен из метаданных
	accessToken, err := p.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, p.log, err, fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization))
	}

	// Получаем UserID из токена
	userID, err := p.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, p.log, err, "failed to fetch user id from token")
	}

	messages, err := p.repo.ChatHistory(ctx, userID, in.InterlocutorID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, p.log, err, "failed to fetch chat messages")
	}

	gprcMessages := make([]*gen.ChatMessage, len(messages))
	for index, message := range messages {
		gprcMessages[index] = message.ConvertToGrpcModel()
	}
	return &gen.ChatHistoryResponse{
		Messages: gprcMessages,
	}, nil
}

func (p *ChatProvider) removeClient(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.clients, id)
}

func (p *ChatProvider) addClient(id string, stream gen.ChatService_ChatServer) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients[id] = stream
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, val := range input {
		if _, ok := seen[val]; !ok {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}

func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}
