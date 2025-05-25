package helpers

import (
	"encoding/json"
	"net/http"
)

// ResponseJSON sends a JSON response with the given status code and data.
func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ResponseSuccess(w http.ResponseWriter, messages ...string) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"messages": messages,
	})
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func ResponseNotFound(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
