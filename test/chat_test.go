package test

import (
	"context"
	"errors"
	"finalProjectGolang/chat"
	"finalProjectGolang/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetChatByID(chatID int64) (*database.Chat, error) {
	args := m.Called(chatID)
	return args.Get(0).(*database.Chat), args.Error(1)
}

func (m *MockDatabase) GetMessages(chatID int64) ([]database.Message, error) {
	args := m.Called(chatID)
	return args.Get(0).([]database.Message), args.Error(1)
}

func (m *MockDatabase) GetChat(receiver, sender string) (int64, error) {
	args := m.Called(receiver, sender)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDatabase) CreateChat(receiver, sender string) (int64, error) {
	args := m.Called(receiver, sender)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDatabase) SendMessage(sender, text string, chatID int64) (int64, error) {
	args := m.Called(sender, text, chatID)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetChatByID(t *testing.T) {
	mockDB := new(MockDatabase)
	server := &chat.Server{DB: mockDB}

	req := &chat.GetChatRequest{Chatid: 1, Username: "user1"}
	expectedChat := &database.Chat{ChatID: 1, Username1: "user1", Username2: "user2"}
	expectedMessages := []database.Message{
		{MessageID: 1, Sender: "user1", Text: "Hello"},
	}

	mockDB.On("GetChatByID", req.Chatid).Return(expectedChat, nil)
	mockDB.On("GetMessages", req.Chatid).Return(expectedMessages, nil)

	resp, err := server.GetChatByID(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Chatid, resp.Chatid)
	assert.Equal(t, expectedChat.Username1, resp.Username1)
	assert.Equal(t, expectedChat.Username2, resp.Username2)
	assert.Len(t, resp.Messages, len(expectedMessages))
	assert.Equal(t, expectedMessages[0].MessageID, resp.Messages[0].MessageID)

	mockDB.AssertExpectations(t)
}

func TestSendMessage(t *testing.T) {
	mockDB := new(MockDatabase)
	server := &chat.Server{DB: mockDB}

	msg := &chat.Message{Sender: "user1", Receiver: "user2", Text: "Hello"}
	expectedChatID := int64(1)
	expectedMessageID := int64(1)

	mockDB.On("GetChat", msg.Receiver, msg.Sender).Return(expectedChatID, nil)
	mockDB.On("SendMessage", msg.Sender, msg.Text, expectedChatID).Return(expectedMessageID, nil)

	resp, err := server.SendMessage(context.Background(), msg)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Message sent, messageID:"+string(expectedMessageID), resp.Status)

	mockDB.AssertExpectations(t)
}

func TestSendMessage_NewChat(t *testing.T) {
	mockDB := new(MockDatabase)
	server := &chat.Server{DB: mockDB}

	msg := &chat.Message{Sender: "user1", Receiver: "user2", Text: "Hello"}
	expectedChatID := int64(1)
	expectedMessageID := int64(1)

	mockDB.On("GetChat", msg.Receiver, msg.Sender).Return(int64(0), errors.New("no rows in result set"))
	mockDB.On("CreateChat", msg.Receiver, msg.Sender).Return(expectedChatID, nil)
	mockDB.On("SendMessage", msg.Sender, msg.Text, expectedChatID).Return(expectedMessageID, nil)

	resp, err := server.SendMessage(context.Background(), msg)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Message sent, messageID:"+string(expectedMessageID), resp.Status)

	mockDB.AssertExpectations(t)
}
