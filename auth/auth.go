package auth

import (
	"context"
	db "finalProjectGolang/database"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var jwtKey = []byte("my_secret_key")

type server struct {
	DB *db.Database
}

func (s *server) mustEmbedUnimplementedAuthServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewAuthServiceServer() *server {
	db, err := db.CreateDatabase()
	if err != nil {
		return nil
	}
	return &server{DB: db}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *server) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Check if user already exists
	fmt.Println("asdad")
	//if _, exists := s.users[req.Username]; exists {
	//	return nil, errors.New("user already exists")
	//}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Store user

	err = s.DB.CreateUser(req.Fullname, req.Username, string(hashedPassword))

	if err != nil {
		return nil, err
	}

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
	//// Check if user exists
	//hashedPassword, exists := s.users[req.Username]
	//if !exists {
	//	return nil, errors.New("invalid username or password")
	//}
	//
	//// Compare password
	//err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	//if err != nil {
	//	return nil, errors.New("invalid username or password")
	//}
	//
	//// Generate token
	//expirationTime := time.Now().Add(24 * time.Hour)
	//claims := &Claims{
	//	Username: req.Username,
	//	StandardClaims: jwt.StandardClaims{
	//		ExpiresAt: expirationTime.Unix(),
	//	},
	//}
	//
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//tokenString, err := token.SignedString(jwtKey)
	//if err != nil {
	//	return nil, err
	//}

	return &LoginResponse{Token: "1"}, nil
}
