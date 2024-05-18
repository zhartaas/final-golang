package models

type User struct {
	UserID   int64  `json:"user_id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
