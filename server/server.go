package server

import (
	"Todo/handlers"
	"Todo/middlewares"
	"Todo/utils"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute  // Read timeout for HTTP server
	readHeaderTimeout = 30 * time.Second // Timeout for reading the headers
	writeTimeout      = 5 * time.Minute  // Write timeout for the HTTP server
)

// SetupRoutes provides all the routes that can be used
func SetupRoutes() *Server {
	router := chi.NewRouter()

	// Use common middlewares for all routes
	router.Use(middlewares.CommonMiddlewares()...)

	// Define routes under the "/v1" path
	router.Route("/v1", func(r chi.Router) {
		// Health check endpoint to verify the server is running
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "server is running"})
		})
		// User registration and login routes
		r.Post("/register", handlers.RegisterUser)
		r.Post("/login", handlers.LoginUser)
	})

	// Protected routes requiring authentication
	router.Group(func(r chi.Router) {
		// Apply authentication middleware to the following routes
		r.Use(middlewares.Authenticate)

		// User profile and logout routes
		r.Get("/v1/profile", handlers.UserProfile)
		r.Post("/v1/logout", handlers.LogoutUser)
		r.Delete("/v1/delete-account", handlers.DeleteUser)

		// Routes related to Todo management
		r.Route("/v1/todos", func(r chi.Router) {
			r.Post("/create", handlers.CreateTodo)           // Create a new todo
			r.Get("/search", handlers.SearchTodo)            // Search todos
			r.Get("/all-todos", handlers.GetAllTodos)        // Get all todos
			r.Get("/incomplete", handlers.IncompleteTodo)    // Get incomplete todos
			r.Get("/completed", handlers.CompletedTodo)      // Get completed todos
			r.Put("/mark-completed", handlers.MarkCompleted) // Mark a todo as completed
			r.Delete("/delete", handlers.DeleteTodo)         // Delete a specific todo
			r.Delete("/delete-all", handlers.DeleteAllTodos) // Delete all todos
		})
	})

	// Return the server instance with the configured router
	return &Server{
		Router: router,
	}
}

// Run runs the server
func (svc *Server) Run(port string) error {
	// Create and configure the HTTP server
	svc.server = &http.Server{
		Addr:              port,              // Server address and port
		Handler:           svc.Router,        // HTTP handler for the server
		ReadTimeout:       readTimeout,       // Set read timeout for requests
		ReadHeaderTimeout: readHeaderTimeout, // Set header read timeout
		WriteTimeout:      writeTimeout,      // Set write timeout for responses
	}
	// Start the server and listen for incoming requests
	return svc.server.ListenAndServe()
}

// Shutdown shuts down the server
func (svc *Server) Shutdown(timeout time.Duration) error {
	// Create a context with a timeout to gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Shut down the server using the context with the timeout
	return svc.server.Shutdown(ctx)
}
