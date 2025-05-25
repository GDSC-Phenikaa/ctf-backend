package routes

import (
	"encoding/json"
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Profile returns a handler that responds with the authenticated user's profile.
func Profile(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var user models.User
		if err := database.First(&user, userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func ProfileRoutes(database *gorm.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.AuthMiddleware)
	r.Get("/", Profile(database)) // Get user profile
	return r
}
