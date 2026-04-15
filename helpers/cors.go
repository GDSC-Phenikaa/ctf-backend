package helpers

import (
	"net/http"
	"strings"

	"github.com/GDSC-Phenikaa/ctf-backend/env"
)

func getAllowedOrigin(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return "*"
	}

	allowed := env.AllowedOrigin()
	if allowed == "" {
		allowed = "http://localhost:3000" // default for local dev
	}

	origins := strings.Split(allowed, ",")
	for _, o := range origins {
		o = strings.TrimSpace(o)
		if origin == o || o == "*" {
			return origin
		}
	}

	return strings.TrimSpace(origins[0])
}

// CORSOptionsHandler handles preflight requests
func CORSOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(r))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, Origin, X-Requested-With, Cookie, Cache-Control")
	w.WriteHeader(http.StatusNoContent)
}

// CORSMiddleware sets CORS headers on all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", getAllowedOrigin(r))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, Origin, X-Requested-With, Cookie, Cache-Control")
		next.ServeHTTP(w, r)
	})
}
