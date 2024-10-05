package handlers

import (
	"Todo/database/dbHelper"
	"Todo/middlewares"
	"Todo/models"
	"Todo/utils"
	"net/http"
)

// RegisterUser handles the user registration process
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body models.User

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "name is required")
		return
	}

	if body.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "email is required")
		return
	}

	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}

	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}

	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}

	saveErr := dbHelper.CreateUser(body.Name, body.Email, hashedPassword)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to save user")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"user created successfully"})
}

// LoginUser handles user login and JWT generation
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.UserLogin

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if body.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "email is required")
		return
	}

	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	userID, name, userErr := dbHelper.GetUserInfo(body.Email, body.Password)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "failed to find user")
		return
	}

	if userID == "" || name == "" {
		utils.RespondError(w, http.StatusNotFound, nil, "user not found")
		return
	}

	sessionID, crtErr := dbHelper.CreateUserSession(userID)
	if crtErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, crtErr, "failed to create user session")
		return
	}

	token, genErr := utils.GenerateJWT(userID, name, body.Email, sessionID)
	if genErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, genErr, "failed to generate token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{"user logged in successfully", token})
}

// UserProfile returns the profile information of the logged-in user
func UserProfile(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	userProfile, getErr := dbHelper.GetUserProfile(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get user profile")
		return
	}

	utils.RespondJSON(w, http.StatusOK, userProfile)
}

// LogoutUser terminates the user's session by deleting the session
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	sessionID := userCtx.SessionID

	saveErr := dbHelper.DeleteUserSession(sessionID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete user session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"user logged out successfully"})
}

// DeleteUser deletes the user's account and session
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID
	sessionID := userCtx.SessionID

	saveErr := dbHelper.DeleteUser(userID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete user account")
		return
	}

	saveErr = dbHelper.DeleteUserSession(sessionID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete user session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"account deleted successfully"})
}
