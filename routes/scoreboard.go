package routes

import (
	"encoding/json"
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type ScoreboardEntry struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Score    int    `json:"score"`
}

func ScoreboardRoutes(database *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Options("/*", helpers.CORSOptionsHandler)
	
	// Public endpoints
	r.Get("/ctf", getCTFScoreboardHandler(database))
	r.Get("/lms", getLMSScoreboardHandler(database))

	return r
}

func getCTFScoreboardHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var scoreboard []ScoreboardEntry

		// Aggregate sum of points from challanges table leveraging the solves table join where correct = true
		err := database.Table("users").
			Select("users.id as user_id, users.username, users.name, COALESCE(SUM(challanges.points), 0) as score").
			Joins("LEFT JOIN solves ON solves.user_id = users.id AND solves.correct = ?", true).
			Joins("LEFT JOIN challanges ON challanges.id = solves.challenge_id").
			Where("users.is_admin = ?", false).
			Group("users.id").
			Order("score DESC").
			Scan(&scoreboard).Error

		if err != nil {
			helpers.Error("Failed to aggregate CTF scoreboard: %v", err)
			http.Error(w, "Failed to retrieve scoreboard", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"scoreboard": scoreboard,
		})
	}
}

func getLMSScoreboardHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var scoreboard []ScoreboardEntry

		// Aggregate sum of points from questions table leveraging the question_solves table join where correct = true
		err := database.Table("users").
			Select("users.id as user_id, users.username, users.name, COALESCE(SUM(questions.points), 0) as score").
			Joins("LEFT JOIN question_solves ON question_solves.user_id = users.id AND question_solves.correct = ?", true).
			Joins("LEFT JOIN questions ON questions.id = question_solves.question_id").
			Where("users.is_admin = ?", false).
			Group("users.id").
			Order("score DESC").
			Scan(&scoreboard).Error

		if err != nil {
			helpers.Error("Failed to aggregate LMS scoreboard: %v", err)
			http.Error(w, "Failed to retrieve scoreboard", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"scoreboard": scoreboard,
		})
	}
}
