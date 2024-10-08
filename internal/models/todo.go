package models

type Todo struct {
	ID        int     `json:"id"`
	UserID    float64 `json:"user_id"` // Associate with a user
	Title     string  `json:"title"`
	Completed bool    `json:"completed"`
}
