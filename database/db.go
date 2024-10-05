package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	Todo *sqlx.DB
)

type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
)

// ConnectAndMigrate function connects with a given database and returns error if there is any error
func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode SSLMode) error {
	// Format the connection string using provided parameters
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)
	// Open a connection to the PostgresSQL database using sqlx
	DB, err := sqlx.Open("postgres", connStr)

	if err != nil {
		// Return error if database connection fails
		return err
	}

	// Ping the database to ensure connection is alive
	err = DB.Ping()
	if err != nil {
		// Return error if ping fails (i.e., connection is not valid)
		return err
	}
	// Assign the connected database to the global variable Todo
	Todo = DB
	// Perform database migration after successful connection
	return migrateUp(DB)
}

// ShutdownDatabase closes the database connection gracefully
func ShutdownDatabase() error {
	// Close the global database connection
	return Todo.Close()
}

// migrateUp function migrates the database and handles the migration logic
func migrateUp(db *sqlx.DB) error {
	// Create a new migration driver for PostgresSQL using the database instance
	driver, driErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driErr != nil {
		// Return error if migration driver creation fails
		return driErr
	}
	// Set up the migration instance with the file path for migrations and the database driver
	m, migErr := migrate.NewWithDatabaseInstance(
		"file://database/migrations", // Path to migration files
		"postgres", driver)           // Database driver and name

	if migErr != nil {
		// Return error if migration instance creation fails
		return migErr
	}
	// Run the migrations (updating the database schema)
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		// If an error occurs, but it's not "ErrNoChange" (no changes detected), return the error
		return err
	}
	// Return nil if migration was successful or no changes were needed
	return nil
}
