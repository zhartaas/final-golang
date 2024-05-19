package main

import (
	"finalProjectGolang/auth"
	"finalProjectGolang/chat"
	handlers "finalProjectGolang/handlers"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type servers struct {
	*handlers.AuthServer
	*handlers.ChatServer
}

func main() {
	// Initialize the authServer
	authServer := auth.NewAuthServiceServer()
	chatServer := chat.NewChatServiceServer()
	if authServer == nil || chatServer == nil {
		panic("authserviceserver error")

	}
	defer authServer.DB.Close()

	// Create the server struct with the authServer
	srv := &servers{
		&handlers.AuthServer{authServer},
		&handlers.ChatServer{ChatServer: chatServer},
	}

	// Запуск gRPC сервера
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		auth.RegisterAuthServiceServer(s, authServer)
		chat.RegisterChatServiceServer(s, chatServer)
		log.Println("gRPC server listening on port 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Запуск HTTP сервера
	http.HandleFunc("/register", srv.Register)
	http.HandleFunc("/login", srv.Login)
	http.HandleFunc("/sendmessage", srv.SendMessage)
	log.Println("HTTP server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
