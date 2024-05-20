package user

import (
	"context"
	"errors"
	db "finalProjectGolang/database"
	"fmt"
)

type server struct {
	DB *db.Database
}

func (s *server) mustEmbedUnimplementedUserServiceServer() {
	//TODO implement me
	panic("implement me")
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

func (s *server) GetProfile(ctx context.Context, req *ProfileRequest) (*ProfileResponse, error) {
	res, err := s.DB.GetProfile(req.Id)
	if err != nil {
		return &ProfileResponse{}, err
	}

	if res == nil {
		return &ProfileResponse{}, errors.New("you dont have access")
	}

	pr := &ProfileResponse{FullName: res.FullName, Username: res.Username}

	fmt.Println(pr)

	return pr, nil
}
