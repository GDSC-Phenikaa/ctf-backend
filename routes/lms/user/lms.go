package lmsuser

import (
	"encoding/json"
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func UserLMSRoutes(database *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Options("/*", helpers.CORSOptionsHandler)

	r.With(middlewares.AuthMiddleware).Get("/modules", listModulesHandler(database))
	r.With(middlewares.AuthMiddleware).Get("/lessons/{id}", getLessonHandler(database))
	r.With(middlewares.AuthMiddleware).Post("/questions/{id}/submit", submitQuestionHandler(database))
	r.With(middlewares.AuthMiddleware).Get("/progress", getProgressHandler(database))

	return r
}

func listModulesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		if err := db.Preload("Lessons").Find(&modules).Error; err != nil {
			http.Error(w, "Failed to retrieve modules", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"modules": modules,
		})
	}
}

func getLessonHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessonID := chi.URLParam(r, "id")
		var lesson models.Lesson
		if err := db.Preload("Questions").First(&lesson, lessonID).Error; err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}

		// Security: hide the correct answer from the user response!
		for i := range lesson.Questions {
			lesson.Questions[i].CorrectAnswer = ""
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"lesson": lesson,
		})
	}
}

func submitQuestionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		questionID := chi.URLParam(r, "id")
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var payload struct {
			Answer string `json:"answer"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var question models.Question
		if err := db.First(&question, questionID).Error; err != nil {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}

		correct := question.CorrectAnswer == payload.Answer

		solve := models.QuestionSolve{
			QuestionID:      question.ID,
			UserID:          userID,
			SubmittedAnswer: payload.Answer,
			Correct:         correct,
		}

		if err := db.Create(&solve).Error; err != nil {
			http.Error(w, "Failed to record answer", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"correct": correct,
		})
	}
}

func getProgressHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var solves []models.QuestionSolve
		if err := db.Where("user_id = ?", userID).Find(&solves).Error; err != nil {
			http.Error(w, "Failed to get progress", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"solves": solves,
		})
	}
}
