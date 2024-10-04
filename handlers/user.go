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

	// Parse the request body into a User struct
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	// Validate required fields
	if body.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "name is required")
		return
	}

	if body.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "email is required")
		return
	}

	// Check if user already exists
	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}

	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}

	// Validate password length
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	// Hash the password before saving
	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}

	// Save user information
	saveErr := dbHelper.CreateUser(body.Name, body.Email, hashedPassword)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to save user")
		return
	}

	// Respond with success message
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"user created successfully"})
}

// LoginUser handles user login and JWT generation
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.UserLogin

	// Parse login credentials from request body
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	// Validate required fields
	if body.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "email is required")
		return
	}

	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	// Verify user credentials
	userID, name, userErr := dbHelper.GetUserInfo(body.Email, body.Password)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "failed to find user")
		return
	}

	// Check if user exists
	if userID == "" || name == "" {
		utils.RespondError(w, http.StatusNotFound, nil, "user not found")
		return
	}

	// Create a new user session
	sessionID, crtErr := dbHelper.CreateUserSession(userID)
	if crtErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, crtErr, "failed to create user session")
		return
	}

	// Generate JWT token for the session
	token, genErr := utils.GenerateJWT(userID, name, body.Email, sessionID)
	if genErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, genErr, "failed to generate token")
		return
	}

	// Respond with the generated JWT token
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{"user logged in successfully", token})
}

// UserProfile returns the profile information of the logged-in user
func UserProfile(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	// Fetch user profile based on user ID
	userProfile, getErr := dbHelper.GetUserProfile(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get user profile")
		return
	}

	// Respond with user profile
	utils.RespondJSON(w, http.StatusOK, userProfile)
}

// LogoutUser terminates the user's session by deleting the session
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	sessionID := userCtx.SessionID

	// Delete the user's session
	saveErr := dbHelper.DeleteUserSession(sessionID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete user session")
		return
	}

	// Respond with logout success message
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"user logged out successfully"})
}

// DeleteUser deletes the user's account and session
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID
	sessionID := userCtx.SessionID

	// Delete the user account
	saveErr := dbHelper.DeleteUser(userID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete user account")
		return
	}

	// Delete the user session
	saveErr = dbHelper.DeleteUserSession(sessionID)
	if saveErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete user session")
		return
	}

	// Respond with account deletion success message
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"account deleted successfully"})
}
