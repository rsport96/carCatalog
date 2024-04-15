// Package app provides the application configuration for the MailProxy application.
// It loads the environment variables and sets up the application configuration.
package app

import (
	"catalog/config"
	"catalog/saturator"
	"database/sql"
	// postgres drive
	_ "github.com/lib/pq"
)

// App represents the application configuration.
type App struct {
	// config holds the application configuration.
	config *config.Config

	// db holds the database connection.
	db *sql.DB

	// saturator holds the saturator instance.
	saturator saturator.Saturator
}

// NewApp creates a new App instance by loading the environment variables.
// It returns the new App instance.
func NewApp() *App {
	// Create a new App instance.
	app := &App{}

	// Load the environment variables.
	app.SetConfig("CTG_")

	db, err := sql.Open("postgres", app.Config().DBUrl)
	if err != nil {
		panic(err)
	}
	app.SetDB(db)

	// Set the saturator instance
	app.SetSaturator(app.Config())

	// Return the new App instance.
	return app
}

// SetConfig loads the environment variables with the given prefix and sets the configuration.
func (app *App) SetConfig(prefix string) {
	// Load and set the environment variables.
	app.config = config.LoadEnvVariables(prefix)
}

// Config returns the application configuration.
func (app *App) Config() *config.Config {
	// Return the application configuration.
	return app.config
}

// SetDB sets the database connection.
func (app *App) SetDB(db *sql.DB) {
	// Set the database connection.
	app.db = db
}

// DB returns the database connection.
func (app *App) DB() *sql.DB {
	// Return the database connection.
	return app.db
}

// SetSaturator sets the saturator instance.
func (app *App) SetSaturator(conf *config.Config) {
	// Set the saturator instance.
	app.saturator = saturator.NewSaturator(conf.SatUrl)
}

// Saturator returns the saturator instance.
func (app *App) Saturator() saturator.Saturator {
	// Return the saturator instance.
	return app.saturator
}
