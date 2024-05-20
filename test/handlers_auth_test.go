package test

import (
	"bytes"
	"context"
	"encoding/json"
	"finalProjectGolang/auth"
	"finalProjectGolang/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockAuthServiceClient is a mock implementation of the AuthServiceClient interface
type MockAuthServiceClient struct {
	mock.Mock
}

func (m *MockAuthServiceClient) Register(ctx context.Context, req *auth.RegisterRequest, opts ...grpc.CallOption) (*auth.RegisterResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*auth.RegisterResponse), args.Error(1)
}

func (m *MockAuthServiceClient) Login(ctx context.Context, req *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*auth.LoginResponse), args.Error(1)
}

func TestRegisterHandler(t *testing.T) {
	mockAuth := new(MockAuthServiceClient)
	server := handlers.AuthServer{AuthServer: mockAuth}

	registerRequest := handlers.RegisterRequest{
		Fullname: "Test User",
		Username: "testuser",
		Password: "password123",
	}
	registerRequestJSON, _ := json.Marshal(registerRequest)

	mockAuth.On("Register", mock.Anything, &auth.RegisterRequest{
		Fullname: "Test User",
		Username: "testuser",
		Password: "password123",
	}).Return(&auth.RegisterResponse{Message: "User registered successfully", Token: "mock_token"}, nil)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(registerRequestJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Register)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])
	assert.Equal(t, "mock_token", response["token"])

	mockAuth.AssertExpectations(t)
}

func TestLoginHandler(t *testing.T) {
	mockAuth := new(MockAuthServiceClient)
	server := handlers.AuthServer{AuthServer: mockAuth}

	loginRequest := handlers.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	loginRequestJSON, _ := json.Marshal(loginRequest)

	mockAuth.On("Login", mock.Anything, &auth.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}).Return(&auth.LoginResponse{Token: "mock_token"}, nil)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(loginRequestJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Login)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "mock_token", response["token"])

	mockAuth.AssertExpectations(t)
}
