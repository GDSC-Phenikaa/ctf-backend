package lmsadmin

import (
	"encoding/json"
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func AdminLMSRoutes(database *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Options("/*", helpers.CORSOptionsHandler)

	protected := r.With(middlewares.AuthMiddleware, middlewares.AdminMiddleware(database))

	// Modules CRUD
	protected.Get("/modules", listModulesHandler(database))
	protected.Get("/modules/{id}", getModuleHandler(database))
	protected.Post("/modules", createModuleHandler(database))
	protected.Put("/modules/{id}", updateModuleHandler(database))
	protected.Delete("/modules/{id}", deleteModuleHandler(database))

	// Lessons CRUD
	protected.Get("/lessons", listLessonsHandler(database))
	protected.Get("/lessons/{id}", getLessonHandler(database))
	protected.Post("/lessons", createLessonHandler(database))
	protected.Put("/lessons/{id}", updateLessonHandler(database))
	protected.Delete("/lessons/{id}", deleteLessonHandler(database))

	// Questions CRUD
	protected.Get("/questions", listQuestionsHandler(database))
	protected.Get("/questions/{id}", getQuestionHandler(database))
	protected.Post("/questions", createQuestionHandler(database))
	protected.Put("/questions/{id}", updateQuestionHandler(database))
	protected.Delete("/questions/{id}", deleteQuestionHandler(database))

	return r
}

/* -------------------------------------------------------------------------- */
/*                               MODULE HANDLERS                              */
/* -------------------------------------------------------------------------- */

func listModulesHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		if err := database.Preload("Lessons.Questions").Find(&modules).Error; err != nil {
			http.Error(w, "Failed to fetch modules", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "modules": modules})
	}
}

func getModuleHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var module models.Module
		if err := database.Preload("Lessons.Questions").First(&module, id).Error; err != nil {
			http.Error(w, "Module not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "module": module})
	}
}

func createModuleHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var module models.Module
		if err := json.NewDecoder(r.Body).Decode(&module); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := database.Create(&module).Error; err != nil {
			http.Error(w, "Failed to create module", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "module": module})
	}
}

func updateModuleHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var module models.Module
		if err := database.First(&module, id).Error; err != nil {
			http.Error(w, "Module not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&module); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		if err := database.Save(&module).Error; err != nil {
			http.Error(w, "Failed to update module", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "module": module})
	}
}

func deleteModuleHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := database.Delete(&models.Module{}, id).Error; err != nil {
			http.Error(w, "Failed to delete module", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Module deleted"})
	}
}

/* -------------------------------------------------------------------------- */
/*                               LESSON HANDLERS                              */
/* -------------------------------------------------------------------------- */

func listLessonsHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var lessons []models.Lesson
		if err := database.Find(&lessons).Error; err != nil {
			http.Error(w, "Failed to fetch lessons", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "lessons": lessons})
	}
}

func getLessonHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var lesson models.Lesson
		if err := database.Preload("Questions").First(&lesson, id).Error; err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "lesson": lesson})
	}
}

func createLessonHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var lesson models.Lesson
		if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := database.Create(&lesson).Error; err != nil {
			http.Error(w, "Failed to create lesson", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "lesson": lesson})
	}
}

func updateLessonHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var lesson models.Lesson
		if err := database.First(&lesson, id).Error; err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		if err := database.Save(&lesson).Error; err != nil {
			http.Error(w, "Failed to update lesson", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "lesson": lesson})
	}
}

func deleteLessonHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := database.Delete(&models.Lesson{}, id).Error; err != nil {
			http.Error(w, "Failed to delete lesson", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Lesson deleted"})
	}
}

/* -------------------------------------------------------------------------- */
/*                              QUESTION HANDLERS                             */
/* -------------------------------------------------------------------------- */

func listQuestionsHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var questions []models.Question
		if err := database.Find(&questions).Error; err != nil {
			http.Error(w, "Failed to fetch questions", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "questions": questions})
	}
}

func getQuestionHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var question models.Question
		if err := database.First(&question, id).Error; err != nil {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "question": question})
	}
}

func createQuestionHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var question models.Question
		if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := database.Create(&question).Error; err != nil {
			http.Error(w, "Failed to create question", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "question": question})
	}
}

func updateQuestionHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var question models.Question
		if err := database.First(&question, id).Error; err != nil {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		if err := database.Save(&question).Error; err != nil {
			http.Error(w, "Failed to update question", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "question": question})
	}
}

func deleteQuestionHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := database.Delete(&models.Question{}, id).Error; err != nil {
			http.Error(w, "Failed to delete question", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Question deleted"})
	}
}
