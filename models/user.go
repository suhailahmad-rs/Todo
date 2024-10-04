package models

type User struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UserCtx struct {
	UserID    string `json:"userId" db:"user_id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	SessionID string `json:"sessionId" db:"session_id"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UserProfile struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}
