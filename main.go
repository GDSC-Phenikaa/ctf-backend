// @title           Twilight CTF API
// @version         1.0
// @description     Twilight CTF API
// @host            localhost:3333
// @BasePath        /
package main

import (
	"net/http"

	"github.com/GDSC-Phenikaa/twilight-ctf/db"
	_ "github.com/GDSC-Phenikaa/twilight-ctf/docs" // swagger docs
	"github.com/GDSC-Phenikaa/twilight-ctf/routes"
	"github.com/GDSC-Phenikaa/twilight-ctf/sessions"
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
	_ = godotenv.Load()
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

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Mount("/user", routes.UserRoutes(database))

	http.ListenAndServe(":3333", r)
}
