package main

import (
	"finalProjectGolang/auth"
	db "finalProjectGolang/database"
	handlers "finalProjectGolang/handlers"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	// Запуск Database
	database, err := db.CreateDatabase("postgres", "9999", "localhost", 5432, "postgres")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.DB.Close()
	// Запуск gRPC сервера
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		authServer := auth.NewAuthServiceServer(database)
		auth.RegisterAuthServiceServer(s, authServer)
		log.Println("gRPC server listening on port 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Запуск HTTP сервера
	http.HandleFunc("/register", handlers.Register)
	//http.HandleFunc("/login", handlers.Login)
	log.Println("HTTP server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
