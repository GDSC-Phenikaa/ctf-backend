package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ListChallengesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := middlewares.GetUserID(r) // If not logged in, userID will be 0

		var challenges []models.Challanges
		if err := db.Where("hidden = ?", false).Find(&challenges).Error; err != nil {
			http.Error(w, "Failed to fetch challenges", http.StatusInternalServerError)
			return
		}

		// Get all solves for this user
		var solves []models.Solves
		if userID != 0 {
			db.Where("user_id = ? AND correct = ?", userID, true).Find(&solves)
		}
		solvedMap := make(map[uint]bool)
		helpers.Debug("solvedMap: %+v\n", solvedMap)
		fmt.Println("userID from context:", userID)
		for _, s := range solves {
			solvedMap[s.ChallengeID] = true
		}

		type ChallengeResponse struct {
			ID          uint   `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Difficulty  string `json:"difficulty"`
			Type        string `json:"type"`
			Points      int    `json:"points"`
			CreatedAt   string `json:"created_at"`
			UpdatedAt   string `json:"updated_at"`
			AuthorID    uint   `json:"author_id"`
			AuthorName  string `json:"author_name"`
			Docker      bool   `json:"docker"`
			Solves      int    `json:"solves"`
			Solved      bool   `json:"solved"`
		}

		resp := make([]ChallengeResponse, 0, len(challenges))
		for _, c := range challenges {
			resp = append(resp, ChallengeResponse{
				ID:          c.ID,
				Title:       c.Title,
				Description: c.Description,
				Difficulty:  c.Difficulty,
				Type:        c.Type,
				Points:      c.Points,
				CreatedAt:   c.CreatedAt,
				UpdatedAt:   c.UpdatedAt,
				AuthorID:    c.AuthorID,
				AuthorName:  c.AuthorName,
				Docker:      c.Docker,
				Solves:      c.Solves,
				Solved:      solvedMap[c.ID],
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"challenges": resp,
		})
	}
}

func SubmitChallengeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var submission struct {
			ChallengeID uint   `json:"challenge_id"`
			Flag        string `json:"flag"`
		}
		if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var challenge models.Challanges
		if err := db.First(&challenge, submission.ChallengeID).Error; err != nil {
			http.Error(w, "Challenge not found", http.StatusNotFound)
			return
		}

		// Get user ID from context (middleware)
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		correct := submission.Flag == challenge.Flag

		// Create a new solve entry
		solve := models.Solves{
			ChallengeID: challenge.ID,
			UserID:      userID,
			Flag:        submission.Flag,
			Correct:     correct,
		}
		if err := db.Create(&solve).Error; err != nil {
			http.Error(w, "Failed to record solve", http.StatusInternalServerError)
			return
		}

		if correct {
			challenge.Solves++
			if err := db.Save(&challenge).Error; err != nil {
				http.Error(w, "Failed to update challenge", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Flag submitted successfully"})
		} else {
			helpers.ResponseError(w, http.StatusBadRequest, "Incorrect flag")
		}
	}
}

func UserChallengesRoutes(db *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Options("/*", helpers.CORSOptionsHandler)
	r.With(middlewares.AuthMiddleware).Get("/list", ListChallengesHandler(db))
	r.With(middlewares.AuthMiddleware).Post("/submit", SubmitChallengeHandler(db))
	return r
}
