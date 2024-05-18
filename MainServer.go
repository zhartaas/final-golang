package main

import (
	"finalProjectGolang/auth"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func main() {
	// Запуск gRPC сервера
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		authServer := auth.NewAuthServiceServer()
		auth.RegisterAuthServiceServer(s, authServer)
		log.Println("gRPC server listening on port 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Запуск HTTP сервера
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/validate", handlers.ValidateToken)
	log.Println("HTTP server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
