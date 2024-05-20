package main

import (
	"finalProjectGolang/auth"
	pb1 "finalProjectGolang/auth"
	"finalProjectGolang/chat"
	pb2 "finalProjectGolang/chat"
	_ "finalProjectGolang/docs"
	handlers "finalProjectGolang/handlers"
	"finalProjectGolang/user"
	pb3 "finalProjectGolang/user"
	"flag"
	httpSwagger "github.com/swaggo/http-swagger"
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
	*handlers.UserServer
}

var (
	grpcAddr = flag.String("addr", "localhost:50051", "the address to connect to")
)

// @title Final Project Golang API
// @version 1.0
// @description This is the API documentation for the Final Project Golang server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
func main() {
	// Initialize the authServer
	// Create the server struct with the authServer

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		auth.RegisterAuthServiceServer(s, auth.NewAuthServiceServer())
		chat.RegisterChatServiceServer(s, chat.NewChatServiceServer())
		user.RegisterUserServiceServer(s, user.NewChatServiceServer())
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
	userClient := pb3.NewUserServiceClient(conn)
	if authClient == nil || chatClient == nil {
		panic("authserviceserver error")
	}
	srv := &servers{
		&handlers.AuthServer{AuthServer: authClient},
		&handlers.ChatServer{ChatServer: chatClient},
		&handlers.UserServer{UserServer: userClient},
	}

	http.HandleFunc("/register", srv.Register)
	http.HandleFunc("/login", srv.Login)
	http.HandleFunc("/sendmessage", srv.SendMessage)
	http.HandleFunc("/getchat", srv.GetChatByID)
	http.HandleFunc("/profile", srv.GetProfile)

	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("HTTP server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
