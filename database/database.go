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
	chatID    int64
	username1 string
	username2 string
}

type Database struct {
	DB *pgxpool.Pool
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
		query, username1, username2, username1, username2).Scan(&chat.chatID, &chat.username1, &chat.username2)

	if err != nil {
		return 0, err
	}

	return chat.chatID, nil
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
