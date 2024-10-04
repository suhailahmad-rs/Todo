package main

import (
	"Todo/database"
	"Todo/server"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Constant defining the timeout duration for server shutdown
const shutDownTimeOut = 10 * time.Second

func main() {
	// Create a channel to listen for OS signals like Interrupt (Ctrl+C) or SIGTERM (termination signal)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Set up the HTTP routes using the server package
	srv := server.SetupRoutes()

	// Connect to the database and run migrations
	if err := database.ConnectAndMigrate(
		"localhost",
		"5433",
		"todo",
		"local",
		"local",
		database.SSLModeDisable); err != nil { // Disabling SSL mode for local environment
		// Panic and log the error if database initialization and migration fail
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	// Log successful migration message
	logrus.Print("migration successful!!")

	// Run the server in a separate goroutine, so it doesn't block
	go func() {
		if err := srv.Run(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	logrus.Print("Server started at :8080")

	// Wait until an OS signal is received on the 'done' channel
	<-done

	logrus.Info("shutting down server")

	// Gracefully close the database connection
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

	// Gracefully shut down the server with a timeout
	if err := srv.Shutdown(shutDownTimeOut); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}
}
