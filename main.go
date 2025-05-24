// @title           GDSC CTF API
// @version         1.0
// @description     GDSC CTF API
// @host            localhost:3333
// @BasePath        /
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GDSC-Phenikaa/ctf-backend/db"
	_ "github.com/GDSC-Phenikaa/ctf-backend/docs" // swagger docs
	"github.com/GDSC-Phenikaa/ctf-backend/env"
	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/GDSC-Phenikaa/ctf-backend/routes"
	"github.com/GDSC-Phenikaa/ctf-backend/sessions"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func initialize() (*gorm.DB, error) {
	// Initialize the database connection
	database, err := db.Connect()
	if err != nil {
		panic(err)
	}

	// Ensure the database is migrated
	if err := database.AutoMigrate(); err != nil {
		panic(err)
	}

	return database, nil
}

func main() {
	helpers.ParseFlags()

	_ = godotenv.Load()
	if _, err := os.Stat(".env"); err != nil {
		helpers.Error("Failed to load .env file: %v", err)
		os.Exit(1)
	} else {
		helpers.Information("Loaded environment variables from .env file")
	}

	helpers.Information("Starting GDSC CTF API on port %s", env.Port())

	sessions.InitSessionStore()
	database, err := initialize()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	if env.IsDebug() {
		helpers.Warning("Debug mode is enabled. Swagger documentation is available.")
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	r.Mount("/user", routes.UserRoutes(database))

	helpers.Information("Database type: %s", env.DbType())
	helpers.Information("Database name: %s", env.DbName())
	if env.DbType() == "postgres" {
		helpers.Information("Database DSN: %s", env.DbDsn())
	}

	helpers.Success("API is running on http://localhost:%s", env.Port())
	if env.IsDebug() {
		helpers.Success("Swagger documentation is available at http://localhost:%s/swagger/index.html", env.Port())
	}

	http.ListenAndServe(fmt.Sprintf(":%s", env.Port()), r)
}
