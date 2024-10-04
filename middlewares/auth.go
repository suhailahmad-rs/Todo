package middlewares

import (
	"Todo/database/dbHelper"
	"Todo/models"
	"Todo/utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

type ContextKeys string

const (
	// Key used to store user information in the context
	userContext ContextKeys = "userContext"
)

// Authenticate checks the token validation
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "authorization header missing") // Respond with 401 if missing
			return
		}

		// Extract Bearer token from Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.RespondError(w, http.StatusUnauthorized, nil, "bearer token missing") // Respond with 401 if Bearer is missing
			return
		}

		// Parse the JWT token
		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token uses the correct signing method (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method") // Invalid signing method error
			}
			// Return the JWT secret key for validation
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// If the token is invalid or there was an error parsing, return 401
		if parseErr != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid token")
			return
		}

		// Extract the claims from the token (userID, sessionID, etc.)
		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims") // Respond with 401 for invalid claims
			return
		}

		// Extract sessionID from claims
		sessionID := claimValues["sessionID"].(string)

		// Check if the session has been archived (logged out)
		archivedAt, err := dbHelper.GetArchivedAt(sessionID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "internal server error") // Respond with 500 for internal errors
			return
		}

		// If session is archived, invalidate the token
		if archivedAt != nil {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token") // Respond with 401 for archived session
			return
		}

		// Create a UserCtx object to store user information
		user := &models.UserCtx{
			UserID:    claimValues["userID"].(string),
			Name:      claimValues["name"].(string),
			Email:     claimValues["email"].(string),
			SessionID: sessionID,
		}

		// Store the user information in the request context
		ctx := context.WithValue(r.Context(), userContext, user)
		r = r.WithContext(ctx)

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// UserContext extracts the user information from the request context
func UserContext(r *http.Request) *models.UserCtx {
	// Retrieve the user context from the request
	if user, ok := r.Context().Value(userContext).(*models.UserCtx); ok {
		return user
	}
	// Return nil if user context is not found
	return nil
}
