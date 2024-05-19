package chat

import (
	"context"
	"errors"
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

func (s *server) GetChatByID(ctx context.Context, req *GetChatRequest) (*Chat, error) {
	res, err := s.DB.GetChatByID(req.Chatid)
	if err != nil {
		return &Chat{}, err
	}

	if res.Username1 != req.Username && res.Username2 != req.Username {
		return &Chat{}, errors.New("you dont have access")
	}

	resp, err := s.DB.GetMessages(res.ChatID)

	messages := make([]*ChatMessages, 0)

	for _, message := range resp {
		messages = append(messages,
			&ChatMessages{MessageID: message.MessageID, Sender: message.Sender, Text: message.Text})
	}

	chat := &Chat{Chatid: res.ChatID, Username1: res.Username1, Username2: res.Username2, Messages: messages}

	fmt.Println(chat)

	return chat, nil
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

	res, err := s.DB.SendMessage(msg.Sender, msg.Text, chatID)

	if err != nil {
		return &SendResponse{Status: "error occured"}, err
	}
	fmt.Printf("Received message from %s: %s\n", msg.Sender, msg.Text)
	return &SendResponse{Status: fmt.Sprintf("Message sent, messageID:%v", res)}, nil
}
