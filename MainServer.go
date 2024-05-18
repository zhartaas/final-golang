package main

import (
	"finalProjectGolang/auth"
	handlers "finalProjectGolang/restful"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	// Initialize the authServer
	authServer := auth.NewAuthServiceServer()
	if authServer == nil {
		panic("authserviceserver error")
	}
	defer authServer.DB.Close()

	// Create the server struct with the authServer
	srv := &handlers.Server{
		AuthServer: authServer,
	}

	// Запуск gRPC сервера
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		auth.RegisterAuthServiceServer(s, authServer)
		log.Println("gRPC server listening on port 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Запуск HTTP сервера
	http.HandleFunc("/register", srv.Register)
	http.HandleFunc("/login", srv.Login)
	//http.HandleFunc("/validate", srv.ValidateTokenHandler)
	log.Println("HTTP server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
