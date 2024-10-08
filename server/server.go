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
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Use(middlewares.CommonMiddlewares()...)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "server is running"})
		})
		r.Post("/register", handlers.RegisterUser)
		r.Post("/login", handlers.LoginUser)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate)

			r.Route("/user", func(r chi.Router) {
				r.Get("/profile", handlers.GetUser)
				r.Post("/logout", handlers.LogoutUser)
				r.Delete("/delete", handlers.DeleteUser)
			})

			r.Route("/todo", func(r chi.Router) {
				r.Post("/create", handlers.CreateTodo)
				r.Get("/search", handlers.GetTodo)
				r.Get("/all-todos", handlers.GetAllTodos)
				r.Get("/incomplete", handlers.IncompleteTodos)
				r.Get("/completed", handlers.CompletedTodo)
				r.Put("/mark-completed", handlers.MarkCompleted)
				r.Delete("/delete", handlers.DeleteTodo)
				r.Delete("/delete-all", handlers.DeleteAllTodos)
			})
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
