package auth

import (
	"context"
	db "finalProjectGolang/database"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type server struct {
	database *db.Database
}

func (s *server) mustEmbedUnimplementedAuthServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewAuthServiceServer(db *db.Database) *server {
	return &server{database: db}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *server) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Check if user already exists

	//}

	// Hash the password

	//}

	// Store user

	// Generate token
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
		return nil, err
	}

	return &RegisterResponse{
		Message: "User registered successfully",
		Token:   tokenString,
	}, nil
}

func (s *server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Check if user exists

	// Compare password

	// Generate token
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
		return nil, err
	}

	return &LoginResponse{Token: tokenString}, nil
}
