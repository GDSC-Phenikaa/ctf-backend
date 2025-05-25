package admin

import (
	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func AdminRoutes(database *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.With(middlewares.AuthMiddleware, middlewares.AdminMiddleware(database)).Get("/challenges/list", listChallengesHandler(database))     // Admin-only route to list challenges
	r.With(middlewares.AuthMiddleware, middlewares.AdminMiddleware(database)).Post("/challenges/create", createChallengeHandler(database)) // Admin-only route to create a challenge
	r.With(middlewares.AuthMiddleware, middlewares.AdminMiddleware(database)).Put("/challenges/{id}", updateChallengeHandler(database))    // Admin-only route to update a challenge
	return r
}
