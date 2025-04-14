package grpc

import (
	"2025_CakeLand_API/internal/models"
	gen "2025_CakeLand_API/internal/pkg/chat/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/chat/repo"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log/slog"
	"sync"
	"time"
)

type ChatProvider struct {
	gen.UnimplementedChatServiceServer
	clients map[string]gen.ChatService_ChatServer
	log     *slog.Logger
	mu      sync.Mutex
	repo    repo.IChatRepository
}

func NewChatProvider(log *slog.Logger, repo repo.IChatRepository) *ChatProvider {
	return &ChatProvider{
		clients: make(map[string]gen.ChatService_ChatServer),
		log:     log,
		repo:    repo,
	}
}

func (p *ChatProvider) Chat(stream gen.ChatService_ChatServer) error {
	var fromID string

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			p.removeClient(fromID)
			return nil
		}
		if err != nil {
			p.log.Warn("Error receiving message from server", "error", err)
			return err
		}

		// Если это первое сообщение, запоминаем клиента
		if fromID == "" {
			fromID = msg.OwnerID
			p.addClient(fromID, stream)
		}

		// Если нет адресата, ничего не делаем
		if msg.ReceiverID == "" {
			continue
		}

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
				OwnerID:      msg.OwnerID,
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
