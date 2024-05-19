package handlers

import (
	"context"
	"encoding/json"
	chat "finalProjectGolang/chat"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
)

type ChatServer struct {
	ChatServer chat.ChatServiceClient
}

type ChatMessageJSON struct {
	MessageID int64  `json:"message_id"`
	Sender    string `json:"sender"`
	Text      string `json:"text"`
}

type ChatJSON struct {
	ChatID    int64             `json:"chat_id"`
	Username1 string            `json:"username1"`
	Username2 string            `json:"username2"`
	Messages  []ChatMessageJSON `json:"messages"`
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

func (s *ChatServer) GetChatByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from the URL query parameters
	idFromUrl := r.URL.Query().Get("id")
	if idFromUrl == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(idFromUrl)

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
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	fmt.Println(claims)

	//if err != nil {
	//	http.Error(w, "Invalid token", http.StatusUnauthorized)
	//	return
	//}

	//if !token.Valid {
	//	http.Error(w, "Invalid token", http.StatusUnauthorized)
	//	return
	//}

	res, err := s.ChatServer.GetChatByID(context.Background(), &chat.GetChatRequest{
		Username: claims.Username,
		Chatid:   int64(id),
	})

	// Here you would include your logic to get the resource by ID
	// For example, if you had a method `GetResourceByID` on `AuthServer`:
	//res, err := s.ChatServer.GetResourceByID(context.Background(), &auth.GetResourceByIDRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	messages := make([]ChatMessageJSON, 0, len(res.Messages))
	for _, message := range res.Messages {
		messages = append(messages, ChatMessageJSON{
			MessageID: message.MessageID,
			Sender:    message.Sender,
			Text:      message.Text,
		})
	}

	// Create the ChatJSON response
	chatResponse := ChatJSON{
		ChatID:    res.Chatid,
		Username1: res.Username1,
		Username2: res.Username2,
		Messages:  messages,
	}

	// Encode and write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatResponse)
}
