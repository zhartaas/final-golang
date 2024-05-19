package handlers

import (
	"context"
	"encoding/json"
	chat "finalProjectGolang/chat"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

type ChatServer struct {
	ChatServer chat.ChatServiceServer
}

type SendMessageRequest struct {
	Receiver string `json:"receiver"`
	Text     string `json:"text"`
}

func (s *ChatServer) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Change to http.MethodPost
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract Bearer token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	tokenStr := parts[1]

	// Parse and validate the token
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	fmt.Println(claims)

	//if err != nil || !token.Valid {
	//	http.Error(w, "Invalid token", http.StatusUnauthorized)
	//	return
	//}

	// Decode the request body
	var req SendMessageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Use the extracted username as the sender
	res, err := s.ChatServer.SendMessage(context.Background(), &chat.Message{
		Sender:   claims.Username,
		Receiver: req.Receiver,
		Text:     req.Text,
	})

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"token": res.Status}
	json.NewEncoder(w).Encode(response)
}

func (s *ChatServer) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from the URL query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	// Extract the Bearer token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}
	tokenString := parts[1]

	// Parse and validate the JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Here you would include your logic to get the resource by ID
	// For example, if you had a method `GetResourceByID` on `AuthServer`:
	//res, err := s.ChatServer.GetResourceByID(context.Background(), &auth.GetResourceByIDRequest{Id: id})
	if err != nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	// Return the resource in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(res)
}
