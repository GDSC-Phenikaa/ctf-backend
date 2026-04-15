package lmsuser

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
		if err := db.Preload("Segments").Preload("Questions").First(&lesson, lessonID).Error; err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}

		// Security: hide the correct answer from the user response!
		for i := range lesson.Questions {
			lesson.Questions[i].CorrectAnswer = ""
			lesson.Questions[i].AnswerKey = ""
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
			Answer interface{} `json:"answer"`
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

		correct, normalizedAnswer, err := gradeAnswer(question, payload.Answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		submittedAnswer := normalizeSubmission(payload.Answer)

		solve := models.QuestionSolve{
			QuestionID:       question.ID,
			UserID:           userID,
			SubmittedAnswer:  submittedAnswer,
			NormalizedAnswer: normalizedAnswer,
			Correct:          correct,
		}

		if err := db.Create(&solve).Error; err != nil {
			http.Error(w, "Failed to record answer", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":            "success",
			"correct":           correct,
			"awarded_points":    map[bool]int{true: question.Points, false: 0}[correct],
			"normalized_answer": normalizedAnswer,
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

func gradeAnswer(question models.Question, submitted interface{}) (bool, string, error) {
	answerKey := parseAnswerKey(question.AnswerKey)
	typeName := strings.ToLower(strings.TrimSpace(question.Type))

	switch typeName {
	case "single_choice", "mcq":
		expected := normalizeText(getString(answerKey, "correct", question.CorrectAnswer))
		got := normalizeText(asString(submitted))
		return got == expected, got, nil
	case "multi_choice":
		expectedList := normalizeSlice(getStringArray(answerKey, "correct", question.CorrectAnswer))
		gotList := normalizeSlice(asStringSlice(submitted))
		return slicesEqual(expectedList, gotList), strings.Join(gotList, ","), nil
	case "true_false":
		expectedBool, ok := asBool(answerKey["correct"])
		if !ok {
			expectedBool, ok = asBool(question.CorrectAnswer)
		}
		if !ok {
			return false, "", httpError("invalid true_false answer key")
		}
		gotBool, ok := asBool(submitted)
		if !ok {
			return false, "", httpError("invalid boolean answer")
		}
		normalized := strconv.FormatBool(gotBool)
		return gotBool == expectedBool, normalized, nil
	case "numeric":
		expected, ok := asFloat(answerKey["value"])
		if !ok {
			expected, ok = asFloat(question.CorrectAnswer)
		}
		if !ok {
			return false, "", httpError("invalid numeric answer key")
		}
		tolerance := 0.0
		if v, has := asFloat(answerKey["tolerance"]); has {
			tolerance = v
		}
		got, ok := asFloat(submitted)
		if !ok {
			return false, "", httpError("invalid numeric answer")
		}
		normalized := strconv.FormatFloat(got, 'f', -1, 64)
		return math.Abs(got-expected) <= tolerance, normalized, nil
	case "short_text", "long_text", "text", "code":
		accepted := normalizeSlice(getStringArray(answerKey, "accepted", question.CorrectAnswer))
		got := normalizeText(asString(submitted))
		for _, option := range accepted {
			if got == option {
				return true, got, nil
			}
		}
		return false, got, nil
	default:
		expected := normalizeText(getString(answerKey, "correct", question.CorrectAnswer))
		got := normalizeText(asString(submitted))
		return got == expected, got, nil
	}
}

func parseAnswerKey(raw string) map[string]interface{} {
	if strings.TrimSpace(raw) == "" {
		return map[string]interface{}{}
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return map[string]interface{}{}
	}
	return m
}

func normalizeSubmission(answer interface{}) string {
	switch value := answer.(type) {
	case string:
		return value
	default:
		payload, _ := json.Marshal(value)
		return string(payload)
	}
}

func normalizeText(input string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(input)), " "))
}

func normalizeSlice(values []string) []string {
	normalized := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		n := normalizeText(value)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		normalized = append(normalized, n)
	}
	sort.Strings(normalized)
	return normalized
}

func getString(answerKey map[string]interface{}, key, fallback string) string {
	if value, ok := answerKey[key]; ok {
		return asString(value)
	}
	return fallback
}

func getStringArray(answerKey map[string]interface{}, key, fallback string) []string {
	if value, ok := answerKey[key]; ok {
		return asStringSlice(value)
	}
	if fallback == "" {
		return nil
	}
	return []string{fallback}
}

func asString(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		payload, _ := json.Marshal(typed)
		return string(payload)
	}
}

func asStringSlice(value interface{}) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			result = append(result, asString(item))
		}
		return result
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil
		}
		parts := strings.Split(typed, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	default:
		return []string{asString(value)}
	}
}

func asBool(value interface{}) (bool, bool) {
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		n := normalizeText(typed)
		switch n {
		case "true", "1", "yes", "y":
			return true, true
		case "false", "0", "no", "n":
			return false, true
		default:
			return false, false
		}
	case float64:
		if typed == 1 {
			return true, true
		}
		if typed == 0 {
			return false, true
		}
		return false, false
	default:
		return false, false
	}
}

func asFloat(value interface{}) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func httpError(message string) error {
	return &gradingError{message: message}
}

type gradingError struct {
	message string
}

func (e *gradingError) Error() string {
	return e.message
}
