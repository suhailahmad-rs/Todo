package dbHelper

import (
	"Todo/database"
	"Todo/models"
)

func IsTodoExists(name, userID string) (bool, error) {
	SQL := `SELECT count(id) > 0 as is_exist
			  FROM todos
			  WHERE name = TRIM($1)     
			    AND user_id = $2        
			    AND archived_at IS NULL`

	var check bool
	chkErr := database.Todo.Get(&check, SQL, name, userID)
	return check, chkErr
}

func CreateTodo(name, description, userID string) error {
	SQL := `INSERT INTO todos (name, description, user_id)
			  VALUES (TRIM($1), TRIM($2), $3)`

	_, crtErr := database.Todo.Exec(SQL, name, description, userID)
	return crtErr
}

func GetTodo(name, userID string) (models.Todo, error) {
	SQL := `SELECT id, user_id, name, description, is_completed
              FROM todos
              WHERE name ILIKE '%' || $1 || '%' 
                AND user_id = $2              
                AND archived_at IS NULL`

	var todo models.Todo
	getErr := database.Todo.Get(&todo, SQL, name, userID)
	return todo, getErr
}

func GetAllTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, is_completed
			  FROM todos
			  WHERE user_id = $1             
			    AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	getErr := database.Todo.Select(&todos, query, userID)
	return todos, getErr
}

func GetIncompleteTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, is_completed
			  FROM todos
			  WHERE user_id = $1             
			    AND is_completed = false     
			    AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	getErr := database.Todo.Select(&todos, query, userID)
	return todos, getErr
}

func GetCompletedTodos(userID string) ([]models.Todo, error) {
	query := `SELECT id, user_id, name, description, is_completed
			  FROM todos
			  WHERE user_id = $1             
			    AND is_completed = true      
			    AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	getErr := database.Todo.Select(&todos, query, userID)
	return todos, getErr
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
