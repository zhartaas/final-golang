package handlers

import (
	"context"
	"encoding/json"
	"finalProjectGolang/user"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
)

type UserServer struct {
	UserServer user.UserServiceClient
}

type ProfileJSON struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
}

func (s *UserServer) GetProfile(w http.ResponseWriter, r *http.Request) {
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

	res, err := s.UserServer.GetProfile(context.Background(), &user.ProfileRequest{
		Id: int64(id),
	})

	// Here you would include your logic to get the resource by ID
	// For example, if you had a method `GetResourceByID` on `AuthServer`:
	//res, err := s.ChatServer.GetResourceByID(context.Background(), &auth.GetResourceByIDRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create the ChatJSON response
	pr := ProfileJSON{
		FullName: res.FullName,
		Username: res.Username,
	}

	// Encode and write the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pr)
}
