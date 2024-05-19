package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	username = "postgres"
	password = "5641"
	hostname = "localhost"
	port     = 5432
	dbname   = "finalproject"
)

type User struct {
	UserID   int64
	FullName string
	Username string
	Password string
}

type Chat struct {
	ChatID    int64
	Username1 string
	Username2 string
}

type Database struct {
	DB *pgxpool.Pool
}

type Message struct {
	MessageID int64
	Sender    string
	Text      string
	chatID    int64
}

var db *Database

func (db Database) SendMessage(username, text string, chatID int64) (int64, error) {
	query := "INSERT INTO message(sender, text, chatID) VALUES ($1, $2, $3) RETURNING messageID"

	var messageID int64

	err := db.DB.QueryRow(context.Background(), query, username, text, chatID).Scan(&messageID)

	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func (db Database) GetMessages(chatID int64) ([]*Message, error) {
	messages := make([]*Message, 0)

	query := "SELECT * FROM message WHERE chatID = $1"

	rows, err := db.DB.Query(context.Background(), query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		message := &Message{}
		err := rows.Scan(&message.MessageID, &message.Sender, &message.Text, &message.chatID)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (db Database) GetChatByID(id int64) (*Chat, error) {
	var chat Chat

	query := "SELECT * FROM chat WHERE chatID = $1"

	err := db.DB.QueryRow(context.Background(), query, id).Scan(&chat.ChatID, &chat.Username1, &chat.Username2)

	if err != nil {
		return &Chat{}, err
	}
	return &chat, nil
}

func (db Database) GetUser(username string) (*User, error) {
	var user User

	err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE username = $1", username).Scan(&user.UserID, &user.FullName, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given ID
			return &User{}, fmt.Errorf("question with username %d not found", username)
		}
		// Other error occurred
		return &User{}, err
	}

	return &user, nil
}

func (db Database) GetChat(username1, username2 string) (int64, error) {
	query := "SELECT * FROM chat WHERE (username1 = $1 OR username1 = $2) AND (username2 = $3 OR username2 = $4)"

	var chat Chat
	err := db.DB.QueryRow(context.Background(),
		query, username1, username2, username1, username2).Scan(&chat.ChatID, &chat.Username1, &chat.Username2)

	if err != nil {
		return 0, err
	}

	return chat.ChatID, nil
}

func (db Database) CreateChat(username1, username2 string) (int64, error) {
	query := "INSERT INTO chat (username1, username2) VALUES ($1, $2) RETURNING chatID"

	var chatID int64

	err := db.DB.QueryRow(context.Background(), query, username1, username2).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (db Database) CreateUser(fullname, username, password string) error {

	query := "INSERT INTO users (fullname, username, password) VALUES ($1, $2, $3)"

	_, err := db.DB.Exec(context.Background(), query, fullname, username, password)

	if err != nil {
		fmt.Println(err, 1)
		panic(err)
	}
	return nil
}

func (db Database) Close() {
	db.Close()
}

func CreateDatabase() (*Database, error) {
	if db != nil {
		return db, nil
	}
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, dbname)
	DSN = "postgres://finalproject_ub3y_user:AoeGlMcdAr3GNGhl81dwFqOR4lrUIRnd@dpg-cp57sf779t8c73eqdeeg-a/finalproject_ub3y"
	DB, err := pgxpool.Connect(context.Background(), DSN)

	if err != nil {

		panic(err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		//fmt.Println(err)
		return &Database{}, err
	}
	db = &Database{DB: DB}
	fmt.Println("Successfully connected to postgres")

	return db, nil
}

//
