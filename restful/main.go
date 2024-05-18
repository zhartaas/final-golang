package main

import (
	"context"
	"encoding/json"
	pb "finalProjectGolang/models"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
	"time"
)

type RegisterRequest struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var grpcaddr = flag.String("addr", "localhost:50051", "the address to connect to")
var conn pb.UserServiceClient
var grpcConn *grpc.ClientConn

func main() {
	flag.Parse()

	// Establish the gRPC connection
	if err := connectToGrpc(); err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer grpcConn.Close() // Ensure the connection is closed when the application exits

	// Define HTTP handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	handleHandlers()

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func connectToGrpc() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	grpcConn, err = grpc.DialContext(ctx, *grpcaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	conn = pb.NewUserServiceClient(grpcConn)
	return nil
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Fullname == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	res, err := conn.CreateUser(context.Background(), &pb.RegisterUserRequest{Fullname: req.Fullname, Username: req.Username, Password: req.Password})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res.Message)

}

func handleGet(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	//jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "get", id)
	//w.Write(jsonResponse)
}

func handleHandlers() {
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/register", handleRegister)
}
