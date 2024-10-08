package handlers

import (
	"Todo/database/dbHelper"
	"Todo/middlewares"
	"Todo/models"
	"Todo/utils"
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var body models.TodoRequest
	userCtx := middlewares.UserContext(r)
	body.UserID = userCtx.UserID

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
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

	if saveErr := dbHelper.CreateTodo(body.Name, body.Description, body.UserID); saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to create todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo created successfully"})
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todo, getErr := dbHelper.GetTodo(name, userID)
	if getErr != nil {
		if errors.Is(getErr, sql.ErrNoRows) {
			utils.RespondError(w, http.StatusOK, getErr, "todo not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get todo")
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, getErr := dbHelper.GetAllTodos(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get todos")
		return
	}

	if len(todos) == 0 {
		utils.RespondError(w, http.StatusOK, getErr, "no todo found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
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
