package routes

import (
	"encoding/json"
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/GDSC-Phenikaa/ctf-backend/sessions"
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
func LoginHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Email string
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var user models.User
		if err := database.Where("email = ?", creds.Email).First(&user).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		sessions.SetUserID(w, r, user.ID)
		w.Write([]byte("Logged in"))
	}
}

// @Summary      Get Profile
// @Description  Gets the current user's profile
// @Tags         user
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      404  {string}  string  "User not found"
// @Router       /user/profile [get]
func ProfileHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := sessions.GetUserID(r)
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

// @Summary      Register
// @Description  Registers a new user
// @Tags         user
// @Accept       json
// @Produce      plain
// @Param        user  body  models.User  true  "User info"
// @Success      201  {string}  string  "Registered"
// @Failure      400  {string}  string  "Bad request"
// @Failure      409  {string}  string  "Email already exists"
// @Router       /user/register [post]
func RegisterHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Check if email already exists
		var existing models.User
		if err := database.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		if err := database.Create(&req).Error; err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Registered"))
	}
}

func UserRoutes(database *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", LoginHandler(database))
	r.Get("/profile", ProfileHandler(database))
	r.Post("/register", RegisterHandler(database))
	return r
}
