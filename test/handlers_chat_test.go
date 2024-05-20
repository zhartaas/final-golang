package test

import (
	"bytes"
	"context"
	"encoding/json"
	"finalProjectGolang/chat"
	"finalProjectGolang/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var jwtKey = []byte("100%forProject")

// MockChatServiceClient is a mock implementation of the ChatServiceClient interface
type MockChatServiceClient struct {
	mock.Mock
}

func (m *MockChatServiceClient) GetChatByID(ctx context.Context, req *chat.GetChatRequest, opts ...grpc.CallOption) (*chat.Chat, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*chat.Chat), args.Error(1)
}

func (m *MockChatServiceClient) SendMessage(ctx context.Context, req *chat.Message, opts ...grpc.CallOption) (*chat.SendResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*chat.SendResponse), args.Error(1)
}

func generateToken(username string) string {
	claims := &handlers.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func TestSendMessageHandler(t *testing.T) {
	mockChat := new(MockChatServiceClient)
	server := handlers.ChatServer{ChatServer: mockChat}

	token := generateToken("user1")

	sendMessageRequest := handlers.SendMessageRequest{
		Receiver: "user2",
		Text:     "Hello",
	}
	sendMessageRequestJSON, _ := json.Marshal(sendMessageRequest)

	mockChat.On("SendMessage", mock.Anything, &chat.Message{
		Sender:   "user1",
		Receiver: "user2",
		Text:     "Hello",
	}).Return(&chat.SendResponse{Status: "Message sent, messageID:1"}, nil)

	req, err := http.NewRequest("POST", "/sendmessage", bytes.NewBuffer(sendMessageRequestJSON))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.SendMessage)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Message sent, messageID:1", response["token"])

	mockChat.AssertExpectations(t)
}

func TestGetChatByIDHandler(t *testing.T) {
	mockChat := new(MockChatServiceClient)
	server := handlers.ChatServer{ChatServer: mockChat}

	token := generateToken("user1")

	expectedChat := &chat.Chat{
		Chatid:    1,
		Username1: "user1",
		Username2: "user2",
		Messages: []*chat.ChatMessages{
			{MessageID: 1, Sender: "user1", Text: "Hello"},
		},
	}

	mockChat.On("GetChatByID", mock.Anything, &chat.GetChatRequest{
		Username: "user1",
		Chatid:   1,
	}).Return(expectedChat, nil)

	req, err := http.NewRequest("GET", "/getchat?id=1", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetChatByID)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response handlers.ChatJSON
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedChat.Chatid, response.ChatID)
	assert.Equal(t, expectedChat.Username1, response.Username1)
	assert.Equal(t, expectedChat.Username2, response.Username2)
	assert.Len(t, response.Messages, 1)
	assert.Equal(t, expectedChat.Messages[0].MessageID, response.Messages[0].MessageID)

	mockChat.AssertExpectations(t)
}
