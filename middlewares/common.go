package middlewares

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
)

// corsOptions setting up routes for CORS
func corsOptions() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                                                                                                                                                 // Allow all origins, adjust as necessary for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},                                                                                                                  // HTTP methods allowed by CORS
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Token", "importDate", "X-Client-Version", "Cache-Control", "Pragma", "x-started-at", "x-api-key"}, // Allowed headers for CORS
		ExposedHeaders:   []string{"Link"},                                                                                                                                                              // Headers that are exposed to the client
		AllowCredentials: true,                                                                                                                                                                          // Allow credentials such as cookies or authorization headers
	})
}

// CommonMiddlewares middleware common for all routes
func CommonMiddlewares() chi.Middlewares {
	return chi.Chain(
		// Middleware to set the Content-Type header to JSON for all responses
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json") // Ensure response is JSON
				next.ServeHTTP(w, r)
			})
		},
		// Apply CORS options to all requests
		corsOptions().Handler,
		// Middleware to recover from panics and log the error
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					err := recover()
					if err != nil {
						logrus.Errorf("Request Panic err: %v", err) // Log the panic error
						// Return a 500 Internal Server Error response in case of panic
						jsonBody, _ := json.Marshal(map[string]string{
							"error": "There was an internal server error",
						})
						w.Header().Set("Content-Type", "application/json") // Set response as JSON
						w.WriteHeader(http.StatusInternalServerError)      // Set status to 500
						_, err := w.Write(jsonBody)                        // Write the error response to the client
						if err != nil {
							logrus.Errorf("Failed to send response from middleware with error: %+v", err) // Log error if response fails to send
						}
					}
				}()
				next.ServeHTTP(w, r) // Proceed with the next handler
			})
		},
	)
}
