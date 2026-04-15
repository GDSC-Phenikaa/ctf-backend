package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/GDSC-Phenikaa/ctf-backend/env"
	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := ""
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			helpers.Information("[AUTH] Found token in Authorization header")
		} else {
			tokenString = r.URL.Query().Get("token")
			if tokenString != "" {
				helpers.Information("[AUTH] Found token in Query String")
			} else {
				// If not in query, check cookie
				if cookie, err := r.Cookie("workspace_token"); err == nil {
					tokenString = cookie.Value
					helpers.Information("[AUTH] Found token in Cookie")
				}
			}
		}

		if tokenString == "" {
			helpers.Warning("[AUTH] No token found in header, query, or cookie")
			http.Error(w, "Missing or invalid Authorization", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				helpers.Warning("[AUTH] Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(env.JwtSecret()), nil
		})

		if err != nil {
			helpers.Warning("[AUTH] JWT Parse Error: %v", err)
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			helpers.Warning("[AUTH] Token is nil or invalid")
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		helpers.Success("[AUTH] Token valid for workspace access")

		// If the token was found in the query string, persist it in a cookie for future requests
		// (needed for JS/CSS/WebSocket assets inside the iframe that won't have the query string)
		if qToken := r.URL.Query().Get("token"); qToken != "" {
			helpers.Information("[AUTH] Setting workspace_token cookie")
			http.SetCookie(w, &http.Cookie{
				Name:     "workspace_token",
				Value:    qToken,
				Path:     "/",
				HttpOnly: true,
				// Secure: true requires HTTPS. For local dev/tunnels, we'll keep it flexible
				Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(24 * time.Hour),
			})
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
