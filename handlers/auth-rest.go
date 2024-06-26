package handlers

import (
	"context"
	"encoding/json"
	auth "finalProjectGolang/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

var jwtKey = []byte("100%forProject")

type AuthServer struct {
	AuthServer auth.AuthServiceClient
}

type RegisterRequest struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// структура для JWT токена
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// @Summary Login a user
// @Description Login a user with username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   username  body  LoginRequest  true  "Username and Password"
// @Success 200 {string} string "ok"
// @Failure 400 {object} string "error"
// @Router /login [post]
func (s *AuthServer) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	res, err := s.AuthServer.Login(context.Background(), &auth.LoginRequest{Username: req.Username, Password: req.Password})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"token": res.Token}
	json.NewEncoder(w).Encode(response)

}

// @Summary Register a new user
// @Description Register a new user with a username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   username  body  string  true  "Username"
// @Param   password  body  string  true  "Password"
// @Success 200 {string} string "ok"
// @Failure 400 {object} string "error"
// @Router /register [post]
func (s *AuthServer) Register(w http.ResponseWriter, r *http.Request) {
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

	// Здесь должна быть логика для сохранения пользователя в базу данных

	res, err := s.AuthServer.Register(context.Background(), &auth.RegisterRequest{Fullname: req.Fullname, Username: req.Username, Password: req.Password})

	if err != nil {
		http.Error(w, "Bad reqeust", http.StatusBadRequest)
		return
	}

	// Возвращаем JWT токен в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": res.Message, "token": res.Token}
	json.NewEncoder(w).Encode(response)
}
