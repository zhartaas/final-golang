package chat

import (
	"context"
	db "finalProjectGolang/database"
	"fmt"
)

type server struct {
	DB *db.Database
}

func (s *server) mustEmbedUnimplementedChatServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewChatServiceServer() *server {
	return &server{}
}

func (s *server) SendMessage(ctx context.Context, msg *Message) (*SendResponse, error) {

	fmt.Printf("Received message from %s: %s\n", msg.Sender, msg.Text)
	return &SendResponse{Status: "Message sent"}, nil
}

func (s *server) ReceiveMessages(req *ReceiveRequest, stream ChatService_ReceiveMessagesServer) error {

	return nil
}
