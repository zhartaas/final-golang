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
type Database struct {
	DB *pgxpool.Pool
}

var db *Database

func (db Database) GetUser(id int) (*User, error) {
	var user User

	err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE userID = $1", id).Scan(&user.UserID, &user.FullName, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given ID
			return &User{}, fmt.Errorf("question with ID %d not found", id)
		}
		// Other error occurred
		return &User{}, err
	}

	return &user, nil
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
