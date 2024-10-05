package dbHelper

import (
	"Todo/database"
	"Todo/models"
	"Todo/utils"
	"database/sql"
	"errors"
	"time"
)

func IsUserExists(email string) (bool, error) {
	query := `SELECT count(id) > 0 as is_exist
			  FROM users
			  WHERE email = TRIM($1)
			    AND archived_at IS NULL`

	var check bool
	chkErr := database.Todo.Get(&check, query, email)
	if chkErr != nil {
		return false, chkErr // Return error if the query fails
	}
	return check, nil
}

// CreateUser Create a new user in the 'users' table.
func CreateUser(name, email, password string) error {
	query := `INSERT INTO users (name, email, password)
			  VALUES (TRIM($1), TRIM($2), $3)`

	_, crtErr := database.Todo.Exec(query, name, email, password)
	if crtErr != nil {
		return crtErr // Return error if the query fails
	}
	return nil
}

// CreateUserSession Create a new user session in the 'user_session' table and return the session ID.
func CreateUserSession(userID string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session(user_id) 
              VALUES ($1) RETURNING id`
	crtErr := database.Todo.QueryRow(query, userID).Scan(&sessionID)

	if crtErr != nil {
		return "", crtErr // Return error if the query fails
	}
	return sessionID, nil
}

// GetUserInfo Retrieve user information by email, along with password validation.
func GetUserInfo(email, password string) (string, string, error) {
	query := `SELECT u.id,
       			   name,
				   u.password
			  FROM users u
			  WHERE u.archived_at IS NULL
			    AND u.email = TRIM($1)`

	var userID string
	var name string
	var passwordHash string
	getErr := database.Todo.QueryRowx(query, email).Scan(&userID, &name, &passwordHash)

	if getErr != nil {
		if errors.Is(getErr, sql.ErrNoRows) {
			return "", "", nil // Return nil if no matching user is found
		}
		return "", "", getErr
	}

	if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
		return "", "", passwordErr // Return error if password validation fails
	}
	return userID, name, nil
}

// GetUserProfile Fetch the user profile for the given user ID.
func GetUserProfile(userID string) (*models.UserProfile, error) {
	var user models.UserProfile
	query := `SELECT id, name, email 
              FROM users 
              WHERE id = $1`

	fetchErr := database.Todo.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if fetchErr != nil {
		if errors.Is(fetchErr, sql.ErrNoRows) {
			return nil, fetchErr // Return error if no matching user is found
		}
		return nil, fetchErr
	}
	return &user, nil
}

// GetArchivedAt Retrieve the 'archived_at' timestamp for a session by session ID.
func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	query := `SELECT archived_at 
              FROM user_session 
              WHERE id = $1`

	getErr := database.Todo.QueryRow(query, sessionID).Scan(&archivedAt)
	if getErr != nil {
		return nil, getErr // Return error if the query fails
	}

	return archivedAt, nil
}

// DeleteUserSession Mark the user session as deleted by updating the 'archived_at' field.
func DeleteUserSession(sessionID string) error {
	query := `UPDATE user_session
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, delErr := database.Todo.Exec(query, sessionID)
	if delErr != nil {
		return delErr
	}
	return nil
}

// DeleteUser Mark the user as deleted by updating the 'archived_at' field.
func DeleteUser(userID string) error {
	query := `UPDATE users
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, delErr := database.Todo.Exec(query, userID)
	if delErr != nil {
		return delErr
	}
	return nil
}
