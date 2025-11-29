package model

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Task struct {
	Id     int64  `json:"id"`
	UserId string `json:"user_id"`
	Title  string `json:"name"`
	Text   string `json:"info"`
}
