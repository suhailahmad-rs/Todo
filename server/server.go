package server

import (
	"Todo/handlers"
	"Todo/middlewares"
	"Todo/utils"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "server is running"})
		})
		r.Post("/register", handlers.RegisterUser)
		r.Post("/login", handlers.LoginUser)
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Authenticate)

		r.Get("/v1/profile", handlers.UserProfile)
		r.Post("/v1/logout", handlers.LogoutUser)
		r.Delete("/v1/delete-account", handlers.DeleteUser)

		r.Route("/v1/todos", func(r chi.Router) {
			r.Post("/create", handlers.CreateTodo)

			// Can use filters in one route using query params
			//start
			r.Get("/search", handlers.SearchTodo)
			r.Get("/all-todos", handlers.GetAllTodos)
			r.Get("/incomplete", handlers.IncompleteTodo)
			r.Get("/completed", handlers.CompletedTodo)
			//end

			r.Put("/mark-completed", handlers.MarkCompleted)
			r.Delete("/delete", handlers.DeleteTodo)
			r.Delete("/delete-all", handlers.DeleteAllTodos)
		})
	})

	return &Server{
		Router: router,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
