package routes

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/GDSC-Phenikaa/ctf-backend/sessions"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterStruct struct {
	Name     string `json:"name" example:"John Doe"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"securepassword"`
}

type LoginStruct struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"securepassword"`
}

// @Summary      Login
// @Description  Logs in a user by email
// @Tags         user
// @Accept       json
// @Produce      plain
// @Param        credentials  body  LoginStruct  true  "User credentials"
// @Success      200  {string}  string  "Logged in"
// @Failure      400  {string}  string  "Bad request"
// @Failure      401  {string}  string  "Invalid credentials"
// @Router       /user/login [post]
func LoginHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			EmailOrUsername string `json:"email"` // Accepts email or username
			Password        string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var user models.User
		if err := database.
			Where("email = ? OR username = ?", creds.EmailOrUsername, creds.EmailOrUsername).
			First(&user).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		// Compare the hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		sessions.SetUserID(w, r, user.ID)

		// Generate JWT token (example)
		token, err := sessions.GenerateJWT(user.ID)
		if err != nil {
			http.Error(w, "Failed to generate token"+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Logged in",
			"token":   token,
		})
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

		type UserProfile struct {
			ID       uint   `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}

		resp := UserProfile{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		}

		json.NewEncoder(w).Encode(resp)
	}
}

// @Summary      Register
// @Description  Registers a new user
// @Tags         user
// @Accept       json
// @Produce      plain
// @Param        user  body  RegisterStruct  true  "User info"
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
		// Username must not contain spaces
		if strings.Contains(req.Username, " ") {
			http.Error(w, "Username must not contain spaces", http.StatusBadRequest)
			return
		}
		// Password requirements: min 8 chars, at least one letter and one number
		var (
			minLen      = 8
			letterRegex = regexp.MustCompile(`[A-Za-z]`)
			numberRegex = regexp.MustCompile(`[0-9]`)
		)
		if len(req.Password) < minLen ||
			!letterRegex.MatchString(req.Password) ||
			!numberRegex.MatchString(req.Password) {
			http.Error(w, "Password must be at least 8 characters and contain at least one letter and one number", http.StatusBadRequest)
			return
		}
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		req.Password = string(hashedPassword)
		if err := database.Create(&req).Error; err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		helpers.ResponseSuccess(w, "Registered")
	}
}

func UserRoutes(database *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", LoginHandler(database))
	r.Get("/profile", ProfileHandler(database))
	r.Post("/register", RegisterHandler(database))
	r.Options("/*", helpers.CORSOptionsHandler)
	return r
}
