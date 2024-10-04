package handlers

import (
	"Todo/database/dbHelper"
	"Todo/middlewares"
	"Todo/models"
	"Todo/utils"
	"net/http"
)

// CreateTodo Handler to create a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var body models.Todo // Struct to hold the request body

	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	body.UserID = userCtx.UserID          // Set UserID in the todo object

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil { // Parse the request body into the 'body' struct
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	// Validate required fields
	if body.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "name is required")
		return
	}

	if body.Description == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "description is required")
		return
	}

	// Check if the todo already exists for the user
	exists, existsErr := dbHelper.IsTodoExists(body.Name, body.UserID)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check todo existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo already exists")
		return
	}

	// Save the new todo in the database
	saveErr := dbHelper.CreateTodo(body.Name, body.Description, body.UserID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to save todo")
		return
	}

	// Respond with success message
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"todo created successfully"})
}

// SearchTodo Handler to search a todo by name
func SearchTodo(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name string `json:"name"`
	}{} // Struct to hold the request body

	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	userID := userCtx.UserID              // Get UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil { // Parse the request body
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.Name == "" { // Validate required field
		utils.RespondError(w, http.StatusBadRequest, nil, "todo name is required")
		return
	}

	// Search for the todo in the database
	todos, getErr := dbHelper.SearchTodo(body.Name, userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to search todo")
		return
	}

	if len(todos) == 0 { // If no todo found, respond with a 404 error
		utils.RespondError(w, http.StatusNotFound, getErr, "no todo found")
		return
	}

	// Respond with the found todos
	utils.RespondJSON(w, http.StatusOK, todos)
}

// GetAllTodos Handler to get all todos for the logged-in user
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	userID := userCtx.UserID              // Get UserID

	todos, getErr := dbHelper.GetAllTodos(userID) // Fetch all todos for the user
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get todos")
		return
	}

	if len(todos) == 0 { // If no todos found, respond with a 404 error
		utils.RespondError(w, http.StatusNotFound, getErr, "no todo found")
		return
	}

	// Respond with the todos
	utils.RespondJSON(w, http.StatusCreated, todos)
}

// IncompleteTodo Handler to get incomplete todos
func IncompleteTodo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	userID := userCtx.UserID              // Get UserID

	todos, getErr := dbHelper.GetIncompleteTodos(userID) // Fetch incomplete todos for the user
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get incomplete todos")
		return
	}

	if len(todos) == 0 { // If no incomplete todos found, respond with a 404 error
		utils.RespondError(w, http.StatusNotFound, getErr, "No todo found")
		return
	}

	// Respond with the incomplete todos
	utils.RespondJSON(w, http.StatusCreated, todos)
}

// CompletedTodo Handler to get completed todos
func CompletedTodo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	userID := userCtx.UserID              // Get UserID

	todos, getErr := dbHelper.GetCompletedTodos(userID) // Fetch completed todos for the user
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get completed todos")
		return
	}

	if len(todos) == 0 { // If no completed todos found, respond with a 404 error
		utils.RespondError(w, http.StatusNotFound, getErr, "No todo found")
		return
	}

	// Respond with the completed todos
	utils.RespondJSON(w, http.StatusCreated, todos)
}

// MarkCompleted Handler to mark a todo as completed
func MarkCompleted(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID string `json:"id"`
	}{} // Struct to hold the request body

	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	userID := userCtx.UserID              // Get UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil { // Parse the request body
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.ID == "" { // Validate required field
		utils.RespondError(w, http.StatusBadRequest, nil, "todo id is required")
		return
	}

	// Mark the todo as completed in the database
	saveErr := dbHelper.MarkCompleted(body.ID, userID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to update todo")
		return
	}

	// Respond with success message
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo marked completed successfully"})
}

// DeleteTodo Handler to delete a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID string `json:"id"`
	}{} // Struct to hold the request body

	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	userID := userCtx.UserID              // Get UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil { // Parse the request body
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.ID == "" { // Validate required field
		utils.RespondError(w, http.StatusBadRequest, nil, "todo id is required")
		return
	}

	// Delete the todo from the database
	saveErr := dbHelper.DeleteTodo(body.ID, userID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete todo")
		return
	}

	// Respond with success message
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo deleted successfully"})
}

// DeleteAllTodos Handler to delete all todos for a user
func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r) // Extract user info from the request context
	userID := userCtx.UserID              // Get UserID

	count, saveErr := dbHelper.DeleteAllTodos(userID) // Delete all todos for the user
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete todos")
		return
	}

	if count == 0 { // If no todos were found to delete, respond with a 404 error
		utils.RespondError(w, http.StatusNotFound, nil, "No todo found")
		return
	}

	// Respond with success message
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"all todos deleted successfully"})
}
