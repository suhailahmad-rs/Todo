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
	var body models.Todo

	userCtx := middlewares.UserContext(r)
	body.UserID = userCtx.UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "name is required")
		return
	}

	if body.Description == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "description is required")
		return
	}

	exists, existsErr := dbHelper.IsTodoExists(body.Name, body.UserID)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check todo existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo already exists")
		return
	}

	saveErr := dbHelper.CreateTodo(body.Name, body.Description, body.UserID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to save todo")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"todo created successfully"})
}

// SearchTodo Handler to search a todo by name
func SearchTodo(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name string `json:"name"`
	}{}

	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo name is required")
		return
	}

	todos, getErr := dbHelper.SearchTodo(body.Name, userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to search todo")
		return
	}

	if len(todos) == 0 {
		utils.RespondError(w, http.StatusNotFound, getErr, "no todo found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

// GetAllTodos Handler to get all todos for the logged-in user
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, getErr := dbHelper.GetAllTodos(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get todos")
		return
	}

	if len(todos) == 0 {
		utils.RespondError(w, http.StatusNotFound, getErr, "no todo found")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, todos)
}

func IncompleteTodo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, getErr := dbHelper.GetIncompleteTodos(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get incomplete todos")
		return
	}

	if len(todos) == 0 {
		utils.RespondError(w, http.StatusNotFound, getErr, "No todo found")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, todos)
}

// CompletedTodo Handler to get completed todos
func CompletedTodo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, getErr := dbHelper.GetCompletedTodos(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get completed todos")
		return
	}

	if len(todos) == 0 {
		utils.RespondError(w, http.StatusNotFound, getErr, "No todo found")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, todos)
}

// MarkCompleted Handler to mark a todo as completed
func MarkCompleted(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID string `json:"id"`
	}{}

	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.ID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo id is required")
		return
	}

	saveErr := dbHelper.MarkCompleted(body.ID, userID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to update todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo marked completed successfully"})
}

// DeleteTodo Handler to delete a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID string `json:"id"`
	}{}

	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.ID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo id is required")
		return
	}

	saveErr := dbHelper.DeleteTodo(body.ID, userID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo deleted successfully"})
}

// DeleteAllTodos Handler to delete all todos for a user
func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	count, saveErr := dbHelper.DeleteAllTodos(userID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete todos")
		return
	}

	if count == 0 {
		utils.RespondError(w, http.StatusNotFound, nil, "No todo found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"all todos deleted successfully"})
}
