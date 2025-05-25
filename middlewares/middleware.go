package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/GDSC-Phenikaa/ctf-backend/env"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(env.JwtSecret()), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if uid, ok := claims["user_id"].(float64); ok {
				ctx := context.WithValue(r.Context(), userIDKey, uint(uid))
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var user struct {
				IsAdmin bool
			}
			if err := db.Table("users").Select("is_admin").Where("id = ?", userID).Scan(&user).Error; err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			if !user.IsAdmin {
				http.Error(w, "Forbidden: admin only", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) (uint, bool) {
	uid, ok := r.Context().Value(userIDKey).(uint)
	return uid, ok
}
