package models

type Todo struct {
	ID          string `json:"id" db:"id"`
	UserID      string `json:"userId" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IsCompleted bool   `json:"isCompleted" db:"is_completed"`
}
