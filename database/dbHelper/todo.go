package dbHelper

import (
	"Todo/database"
	"Todo/models"
)

// IsTodoExists checks if a todo with a given name exists for the specified user.
func IsTodoExists(name, userID string) (bool, error) {
	query := `SELECT count(id) > 0 as is_exist
			  FROM todos
			  WHERE name = TRIM($1)     
			    AND user_id = $2        
			    AND archived_at IS NULL`

	var check bool
	checkErr := database.Todo.Get(&check, query, name, userID)
	if checkErr != nil {
		return false, checkErr
	}
	return check, nil
}

// CreateTodo inserts a new todo into the database for the specified user.
func CreateTodo(name, description, userID string) error {
	query := `INSERT INTO todos (name, description, user_id)
			  VALUES (TRIM($1), TRIM($2), $3)`

	_, createErr := database.Todo.Exec(query, name, description, userID)
	if createErr != nil {
		return createErr
	}
	return nil
}

// SearchTodo searches for todos by name for a specific user.
func SearchTodo(name, userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, is_completed
              FROM todos
              WHERE name ILIKE '%' || $1 || '%' 
                AND user_id = $2              
                AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	searchErr := database.Todo.Select(&todos, query, name, userID)
	return todos, searchErr
}

// GetAllTodos fetches all active todos for the specified user.
func GetAllTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, is_completed
			  FROM todos
			  WHERE user_id = $1             
			    AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	FetchErr := database.Todo.Select(&todos, query, userID)
	return todos, FetchErr
}

// GetIncompleteTodos fetches all incomplete todos for the specified user.
func GetIncompleteTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, is_completed
			  FROM todos
			  WHERE user_id = $1             
			    AND is_completed = false     
			    AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	FetchErr := database.Todo.Select(&todos, query, userID)
	return todos, FetchErr
}

// GetCompletedTodos fetches all completed todos for the specified user.
func GetCompletedTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, is_completed
			  FROM todos
			  WHERE user_id = $1             
			    AND is_completed = true      
			    AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	FetchErr := database.Todo.Select(&todos, query, userID)
	return todos, FetchErr
}

// MarkCompleted marks a specific todo as completed.
func MarkCompleted(id, userID string) error {
	query := `UPDATE todos
              SET is_completed = true        
              WHERE id = $1                  
                AND user_id = $2             
                AND archived_at IS NULL`

	_, updErr := database.Todo.Exec(query, id, userID)
	if updErr != nil {
		return updErr
	}
	return nil
}

// DeleteTodo performs a soft delete by archiving a specific todo.
func DeleteTodo(id, userID string) error {
	query := `UPDATE todos
			  SET archived_at = NOW()        
			  WHERE id = $1                  
			    AND user_id = $2             
			    AND archived_at IS NULL`

	_, delErr := database.Todo.Exec(query, id, userID)
	if delErr != nil {
		return delErr
	}
	return nil
}

// DeleteAllTodos performs a soft delete by archiving all todos for a specific user.
func DeleteAllTodos(userID string) (int, error) {
	query := `UPDATE todos
              SET archived_at = NOW()        
              WHERE user_id = $1             
                AND archived_at IS NULL`

	result, delErr := database.Todo.Exec(query, userID)
	if delErr != nil {
		return 0, delErr
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return 0, rowsErr
	}

	return int(rowsAffected), nil
}
