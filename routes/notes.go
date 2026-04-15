package routes

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type noteResponse struct {
	ID        uint      `json:"id"`
	Href      string    `json:"href"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NotesRoutes(database *gorm.DB) chi.Router {
	r := chi.NewRouter()
	r.Options("/*", helpers.CORSOptionsHandler)
	r.With(middlewares.AuthMiddleware).Get("/", listNotesHandler(database))
	r.With(middlewares.AuthMiddleware).Post("/", createNoteHandler(database))
	r.With(middlewares.AuthMiddleware).Get("/{id}", getNoteHandler(database))
	r.With(middlewares.AuthMiddleware).Put("/{id}", updateNoteHandler(database))
	r.With(middlewares.AuthMiddleware).Delete("/{id}", deleteNoteHandler(database))
	return r
}

func listNotesHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		href := r.URL.Query().Get("href")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 20
		}

		query := database.Model(&models.Note{}).Where("user_id = ?", userID)
		if href != "" {
			query = query.Where("href = ?", href)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
			return
		}

		offset := (page - 1) * limit
		var notes []models.Note
		if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
			http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
			return
		}

		totalPages := 0
		if total > 0 {
			totalPages = int(math.Ceil(float64(total) / float64(limit)))
		}

		helpers.ResponseJSON(w, http.StatusOK, map[string]interface{}{
			"status":      "success",
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"notes":       toNoteResponses(notes),
		})
	}
}

func getNoteHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var note models.Note
		if err := database.Where("id = ? AND user_id = ?", chi.URLParam(r, "id"), userID).First(&note).Error; err != nil {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		helpers.ResponseJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"note":   toNoteResponse(note),
		})
	}
}

func createNoteHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var payload struct {
			Href    string `json:"href"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if payload.Href == "" || payload.Content == "" {
			http.Error(w, "href and content are required", http.StatusBadRequest)
			return
		}

		note := models.Note{UserID: userID, Href: payload.Href, Content: payload.Content}
		if err := database.Create(&note).Error; err != nil {
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}

		helpers.ResponseJSON(w, http.StatusCreated, map[string]interface{}{
			"status": "success",
			"note":   toNoteResponse(note),
		})
	}
}

func updateNoteHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var note models.Note
		if err := database.Where("id = ? AND user_id = ?", chi.URLParam(r, "id"), userID).First(&note).Error; err != nil {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		var payload struct {
			Href    string `json:"href"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if payload.Href == "" || payload.Content == "" {
			http.Error(w, "href and content are required", http.StatusBadRequest)
			return
		}

		note.Href = payload.Href
		note.Content = payload.Content
		if err := database.Save(&note).Error; err != nil {
			http.Error(w, "Failed to update note", http.StatusInternalServerError)
			return
		}

		helpers.ResponseJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"note":   toNoteResponse(note),
		})
	}
}

func deleteNoteHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		result := database.Where("id = ? AND user_id = ?", chi.URLParam(r, "id"), userID).Delete(&models.Note{})
		if result.Error != nil {
			http.Error(w, "Failed to delete note", http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		helpers.ResponseJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Note deleted",
		})
	}
}

func toNoteResponse(note models.Note) noteResponse {
	return noteResponse{
		ID:        note.ID,
		Href:      note.Href,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func toNoteResponses(notes []models.Note) []noteResponse {
	responses := make([]noteResponse, 0, len(notes))
	for _, note := range notes {
		responses = append(responses, toNoteResponse(note))
	}
	return responses
}
