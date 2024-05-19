package main

import (
	"finalProjectGolang/auth"
	pb1 "finalProjectGolang/auth"
	"finalProjectGolang/chat"
	pb2 "finalProjectGolang/chat"
	handlers "finalProjectGolang/handlers"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"time"
)

type servers struct {
	*handlers.AuthServer
	*handlers.ChatServer
}

var (
	grpcAddr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	// Initialize the authServer
	// Create the server struct with the authServer

	// Запуск gRPC сервера
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		auth.RegisterAuthServiceServer(s, auth.NewAuthServiceServer())
		chat.RegisterChatServiceServer(s, chat.NewChatServiceServer())
		log.Println("gRPC server listening on port 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	time.Sleep(time.Second)

	flag.Parse()

	conn, err := grpc.Dial(*grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	authClient := pb1.NewAuthServiceClient(conn)
	chatClient := pb2.NewChatServiceClient(conn)
	if authClient == nil || chatClient == nil {
		panic("authserviceserver error")

	}
	srv := &servers{
		&handlers.AuthServer{AuthServer: authClient},
		&handlers.ChatServer{ChatServer: chatClient},
	}
	// Запуск HTTP сервера
	http.HandleFunc("/register", srv.Register)
	http.HandleFunc("/login", srv.Login)
	http.HandleFunc("/sendmessage", srv.SendMessage)
	http.HandleFunc("/getchat", srv.GetChatByID)
	log.Println("HTTP server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
