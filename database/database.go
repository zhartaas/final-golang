package db

import (
	"context"
	"finalProjectGolang/models"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	DB *pgxpool.Pool
}

func (db *Database) GetAllUsers() ([]models.User, error) {
	query := "SELECT * FROM users"
	var users []models.User

	rows, err := db.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.FullName, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *Database) GetUser(id int) (*models.User, error) {
	var user models.User

	err := db.DB.QueryRow(context.Background(),
		"SELECT * FROM users WHERE userID = $1", id).Scan(&user.UserID, &user.FullName, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) CreateUser(fullname, username, password string) error {
	query := "INSERT INTO users (fullname, username, password) VALUES ($1, $2, $3)"

	_, err := db.DB.Exec(context.Background(), query, fullname, username, password)
	if err != nil {
		return err
	}
	return nil
}

func CreateDatabase(username, password, hostname string, port int, dbname string) (*Database, error) {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, dbname)
	DB, err := pgxpool.Connect(context.Background(), DSN)
	if err != nil {
		return nil, err
	}

	db := &Database{DB: DB}
	fmt.Println("Successfully connected to postgres")
	return db, nil
}
