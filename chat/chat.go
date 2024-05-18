package chat

import (
	"context"
	"fmt"
	"sync"
)

type server struct {
	UnimplementedChatServiceServer
	mu       sync.Mutex
	messages []*Message
}

func NewChatServiceServer() *server {
	return &server{}
}

func (s *server) SendMessage(ctx context.Context, msg *Message) (*SendResponse, error) {
	s.mu.Lock()
	s.messages = append(s.messages, msg)
	s.mu.Unlock()

	fmt.Printf("Received message from %s: %s\n", msg.Sender, msg.Text)
	return &SendResponse{Status: "Message sent"}, nil
}

func (s *server) ReceiveMessages(req *ReceiveRequest, stream ChatService_ReceiveMessagesServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, msg := range s.messages {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
