package dbHelper

import (
	"Todo/database"
	"Todo/models"
	"Todo/utils"
	"database/sql"
	"errors"
	"time"
)

// IsUserExists Check if a user exists in the 'users' table by email.
func IsUserExists(email string) (bool, error) {
	query := `SELECT count(id) > 0 as is_exist
			  FROM users
			  WHERE email = TRIM($1)
			    AND archived_at IS NULL` // Only check for non-archived users by ensuring archived_at is NULL

	var check bool
	chkErr := database.Todo.Get(&check, query, email) // Execute the query and store the result in 'check'
	if chkErr != nil {
		return false, chkErr // Return error if the query fails
	}
	return check, nil
}

// CreateUser Create a new user in the 'users' table.
func CreateUser(name, email, password string) error {
	query := `INSERT INTO users (name, email, password)
			  VALUES (TRIM($1), TRIM($2), $3)` // Insert user data, trimming email and name for consistency

	_, crtErr := database.Todo.Exec(query, name, email, password) // Execute the query to insert user data
	if crtErr != nil {
		return crtErr // Return error if the query fails
	}
	return nil
}

// CreateUserSession Create a new user session in the 'user_session' table and return the session ID.
func CreateUserSession(userID string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session(user_id) 
              VALUES ($1) RETURNING id` // Insert a new session for the user and return the session ID
	crtErr := database.Todo.QueryRow(query, userID).Scan(&sessionID) // Execute the query and store the session ID

	if crtErr != nil {
		return "", crtErr // Return error if the query fails
	}
	return sessionID, nil // Return the session ID
}

// GetUserInfo Retrieve user information by email, along with password validation.
func GetUserInfo(email, password string) (string, string, error) {
	query := `SELECT u.id,
       			   name,
				   u.password
			  FROM users u
			  WHERE u.archived_at IS NULL
			    AND u.email = TRIM($1)` // Select user details only if they are not archived

	var userID string
	var name string
	var passwordHash string
	getErr := database.Todo.QueryRowx(query, email).Scan(&userID, &name, &passwordHash) // Execute the query and scan results

	if getErr != nil {
		if errors.Is(getErr, sql.ErrNoRows) {
			return "", "", nil // Return nil if no matching user is found
		}
		return "", "", getErr // Return error if query fails
	}

	if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
		return "", "", passwordErr // Return error if password validation fails
	}
	return userID, name, nil // Return user ID and name if password is correct
}

// GetUserProfile Fetch the user profile for the given user ID.
func GetUserProfile(userID string) (*models.UserProfile, error) {
	var user models.UserProfile
	query := `SELECT id, name, email 
              FROM users 
              WHERE id = $1` // Select the user's profile details using their ID

	fetchErr := database.Todo.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email) // Execute query and scan result into UserProfile struct
	if fetchErr != nil {
		if errors.Is(fetchErr, sql.ErrNoRows) {
			return nil, fetchErr // Return error if no matching user is found
		}
		return nil, fetchErr // Return error if the query fails
	}
	return &user, nil // Return the user's profile information
}

// GetArchivedAt Retrieve the 'archived_at' timestamp for a session by session ID.
func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	query := `SELECT archived_at 
              FROM user_session 
              WHERE id = $1` // Select the 'archived_at' value from the user session

	getErr := database.Todo.QueryRow(query, sessionID).Scan(&archivedAt) // Execute query and scan the timestamp
	if getErr != nil {
		return nil, getErr // Return error if the query fails
	}

	return archivedAt, nil // Return the 'archived_at' value
}

// DeleteUserSession Mark the user session as deleted by updating the 'archived_at' field.
func DeleteUserSession(sessionID string) error {
	query := `UPDATE user_session
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL` // Update the 'archived_at' field to mark the session as deleted only if it's not already archived

	_, delErr := database.Todo.Exec(query, sessionID) // Execute the query to update the session
	if delErr != nil {
		return delErr // Return error if the update fails
	}
	return nil
}

// DeleteUser Mark the user as deleted by updating the 'archived_at' field.
func DeleteUser(userID string) error {
	query := `UPDATE users
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL` // Update the 'archived_at' field to mark the user as deleted only if it's not already archived

	_, delErr := database.Todo.Exec(query, userID) // Execute the query to update the user record
	if delErr != nil {
		return delErr // Return error if the update fails
	}
	return nil
}
