package main

import (
	"context"
	db "finalProjectGolang/database"
	pb "finalProjectGolang/models"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = flag.Int("port", 50051, "Server port")
var DB *db.Database

type server struct {
	pb.UnimplementedUserServiceServer
}

func main() {
	db, err := db.CreateDatabase()
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *server) CreateUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	err := DB.CreateUser(in.Fullname, in.Username, in.Password)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterUserResponse{Message: "User registered successfully"}, nil
}
