package chat

import (
	"context"
	db "finalProjectGolang/database"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type server struct {
	DB *db.Database
}

func (s *server) mustEmbedUnimplementedChatServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewChatServiceServer() *server {
	db, err := db.CreateDatabase()
	if err != nil {
		return nil
	}
	return &server{DB: db}
}

func (s *server) SendMessage(ctx context.Context, msg *Message) (*SendResponse, error) {
	fmt.Println(s)
	chatID, err := s.DB.GetChat(msg.Receiver, msg.Sender)

	if err != nil {
		fmt.Println(err, pgx.ErrNoRows)
		if err != pgx.ErrNoRows {
			return &SendResponse{Status: "Error occured"}, err
		}
		chatID, err = s.DB.CreateChat(msg.Receiver, msg.Sender)

		if err != nil {
			return &SendResponse{Status: "error occured"}, err
		}
	}

	res, err := s.DB.SendMessage(msg.Receiver, msg.Text, chatID)

	if err != nil {
		return &SendResponse{Status: "error occured"}, err
	}
	fmt.Printf("Received message from %s: %s\n", msg.Sender, msg.Text)
	return &SendResponse{Status: fmt.Sprintf("Message sent, messageID:%v", res)}, nil
}

func (s *server) ReceiveMessages(req *ReceiveRequest, stream ChatService_ReceiveMessagesServer) error {

	return nil
}
