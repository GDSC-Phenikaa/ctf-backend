package lmsadmin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	// Video Segments CRUD
	protected.Get("/video-segments", listVideoSegmentsHandler(database))
	protected.Get("/video-segments/{id}", getVideoSegmentHandler(database))
	protected.Post("/video-segments", createVideoSegmentHandler(database))
	protected.Put("/video-segments/{id}", updateVideoSegmentHandler(database))
	protected.Delete("/video-segments/{id}", deleteVideoSegmentHandler(database))

	// Questions CRUD
	protected.Get("/questions", listQuestionsHandler(database))
	protected.Get("/questions/{id}", getQuestionHandler(database))
	protected.Post("/questions", createQuestionHandler(database))
	protected.Put("/questions/{id}", updateQuestionHandler(database))
	protected.Delete("/questions/{id}", deleteQuestionHandler(database))

	// One-shot migration endpoint for legacy LMS content.
	protected.Post("/migrations/v2", runLMSV2MigrationHandler(database))

	return r
}

/* -------------------------------------------------------------------------- */
/*                               MODULE HANDLERS                              */
/* -------------------------------------------------------------------------- */

func listModulesHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		if err := database.Preload("Lessons.Segments").Preload("Lessons.Questions").Find(&modules).Error; err != nil {
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
		if err := database.Preload("Lessons.Segments").Preload("Lessons.Questions").First(&module, id).Error; err != nil {
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
		if err := database.Preload("Segments").Preload("Questions").Find(&lessons).Error; err != nil {
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
		if err := database.Preload("Segments").Preload("Questions").First(&lesson, id).Error; err != nil {
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
		if lesson.Body == "" {
			lesson.Body = lesson.Content
		}
		if lesson.VideoIframe != "" && !isValidYouTubeIframe(lesson.VideoIframe) {
			http.Error(w, "video_iframe must be a YouTube iframe embed", http.StatusBadRequest)
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
		if lesson.Body == "" {
			lesson.Body = lesson.Content
		}
		if lesson.VideoIframe != "" && !isValidYouTubeIframe(lesson.VideoIframe) {
			http.Error(w, "video_iframe must be a YouTube iframe embed", http.StatusBadRequest)
			return
		}
		if err := database.Save(&lesson).Error; err != nil {
			http.Error(w, "Failed to update lesson", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "lesson": lesson})
	}
}

/* -------------------------------------------------------------------------- */
/*                          VIDEO SEGMENT HANDLERS                            */
/* -------------------------------------------------------------------------- */

func listVideoSegmentsHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var segments []models.VideoSegment
		if err := database.Find(&segments).Error; err != nil {
			http.Error(w, "Failed to fetch video segments", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "video_segments": segments})
	}
}

func getVideoSegmentHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var segment models.VideoSegment
		if err := database.First(&segment, id).Error; err != nil {
			http.Error(w, "Video segment not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "video_segment": segment})
	}
}

func createVideoSegmentHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var segment models.VideoSegment
		if err := json.NewDecoder(r.Body).Decode(&segment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if segment.EndSeconds <= segment.StartSeconds {
			http.Error(w, "end_seconds must be greater than start_seconds", http.StatusBadRequest)
			return
		}
		if err := database.Create(&segment).Error; err != nil {
			http.Error(w, "Failed to create video segment", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "video_segment": segment})
	}
}

func updateVideoSegmentHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var segment models.VideoSegment
		if err := database.First(&segment, id).Error; err != nil {
			http.Error(w, "Video segment not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&segment); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		if segment.EndSeconds <= segment.StartSeconds {
			http.Error(w, "end_seconds must be greater than start_seconds", http.StatusBadRequest)
			return
		}
		if err := database.Save(&segment).Error; err != nil {
			http.Error(w, "Failed to update video segment", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "video_segment": segment})
	}
}

func deleteVideoSegmentHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := database.Delete(&models.VideoSegment{}, id).Error; err != nil {
			http.Error(w, "Failed to delete video segment", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Video segment deleted"})
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
		normalizeQuestionForCompatibility(&question)
		if err := validateQuestionPlacement(database, &question); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		normalizeQuestionForCompatibility(&question)
		if err := validateQuestionPlacement(database, &question); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := database.Save(&question).Error; err != nil {
			http.Error(w, "Failed to update question", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "question": question})
	}
}

/* -------------------------------------------------------------------------- */
/*                               MIGRATIONS                                   */
/* -------------------------------------------------------------------------- */

func runLMSV2MigrationHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx := database.Begin()
		if tx.Error != nil {
			http.Error(w, "Failed to start migration transaction", http.StatusInternalServerError)
			return
		}

		if err := tx.Model(&models.Lesson{}).
			Where("body IS NULL OR body = ''").
			Update("body", gorm.Expr("content")).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to backfill lesson body", http.StatusInternalServerError)
			return
		}

		var questions []models.Question
		if err := tx.Find(&questions).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to load questions", http.StatusInternalServerError)
			return
		}

		migrated := 0
		for i := range questions {
			q := questions[i]

			if q.Prompt == "" {
				q.Prompt = q.Content
			}
			if q.Placement == "" {
				q.Placement = "lesson"
			}
			if q.AnswerKey == "" {
				q.AnswerKey = buildLegacyAnswerKey(q)
			}

			if err := tx.Save(&q).Error; err != nil {
				tx.Rollback()
				http.Error(w, "Failed to migrate questions", http.StatusInternalServerError)
				return
			}
			migrated++
		}

		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Failed to commit migration", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":             "success",
			"migrated_questions": migrated,
		})
	}
}

func normalizeQuestionForCompatibility(question *models.Question) {
	if question.Prompt == "" {
		question.Prompt = question.Content
	}
	if question.Content == "" {
		question.Content = question.Prompt
	}
	if question.Placement == "" {
		question.Placement = "lesson"
	}
	if question.AnswerKey == "" {
		question.AnswerKey = buildLegacyAnswerKey(*question)
	}
}

func validateQuestionPlacement(database *gorm.DB, question *models.Question) error {
	if question.Placement != "lesson" && question.Placement != "segment" {
		return fmt.Errorf("placement must be either lesson or segment")
	}

	if question.Placement == "lesson" {
		question.VideoSegmentID = nil
		return nil
	}

	if question.VideoSegmentID == nil {
		return fmt.Errorf("video_segment_id is required when placement is segment")
	}

	var segment models.VideoSegment
	if err := database.First(&segment, *question.VideoSegmentID).Error; err != nil {
		return fmt.Errorf("video segment not found")
	}
	if segment.LessonID != question.LessonID {
		return fmt.Errorf("video_segment_id must belong to the same lesson")
	}

	return nil
}

func buildLegacyAnswerKey(question models.Question) string {
	legacyType := strings.ToLower(strings.TrimSpace(question.Type))
	correct := strings.TrimSpace(question.CorrectAnswer)

	switch legacyType {
	case "multi_choice":
		parts := strings.Split(correct, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		payload, _ := json.Marshal(map[string]interface{}{"correct": parts})
		return string(payload)
	case "true_false":
		payload, _ := json.Marshal(map[string]interface{}{"correct": strings.EqualFold(correct, "true")})
		return string(payload)
	case "numeric":
		payload, _ := json.Marshal(map[string]interface{}{"value": correct, "tolerance": 0})
		return string(payload)
	case "short_text", "long_text", "text", "code":
		payload, _ := json.Marshal(map[string]interface{}{"accepted": []string{correct}})
		return string(payload)
	default:
		payload, _ := json.Marshal(map[string]interface{}{"correct": correct})
		return string(payload)
	}
}

func isValidYouTubeIframe(raw string) bool {
	normalized := strings.ToLower(raw)
	if !strings.Contains(normalized, "<iframe") {
		return false
	}
	return strings.Contains(normalized, "youtube.com/embed/") || strings.Contains(normalized, "youtube-nocookie.com/embed/")
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
