package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("100%forProject")

type RegisterRequest struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims структура для JWT токена
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}

	http.HandleFunc("/", handler)
	handleHandlers()

	log.Fatal(http.ListenAndServe(":8000", nil))
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

	// Здесь должна быть логика для сохранения пользователя в базу данных

	// Создайте JWT токен
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: req.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Возвращаем JWT токен в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "User registered successfully", "token": tokenString}
	json.NewEncoder(w).Encode(response)
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
