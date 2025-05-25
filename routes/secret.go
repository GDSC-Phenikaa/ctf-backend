package routes

import (
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/env"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// @Summary      Login
// @Description  Logs in a user by email
// @Tags         user
// @Accept       json
// @Produce      plain
// @Param        credentials  body  object  true  "User credentials"
// @Success      200  {string}  string  "Logged in"
// @Failure      400  {string}  string  "Bad request"
// @Failure      401  {string}  string  "Invalid credentials"
// @Router       /user/login [post]
func SecretFlagHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("shhh, don't tell anybody: PKA{" + env.SecretFlag() + "}"))
	}
}

func SecretRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", SecretFlagHandler(nil)) // Pass nil for database if not needed
	return r
}
