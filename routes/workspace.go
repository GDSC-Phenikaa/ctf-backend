package routes

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/GDSC-Phenikaa/ctf-backend/helpers"
	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// StartWorkspace creates a new workspace for the authenticated user
func StartWorkspace(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if user already has an active workspace
		var existing models.Workspace
		if err := db.Where("user_id = ?", userID).First(&existing).Error; err == nil {
			// They already have one running
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(existing)
			return
		}

		// Basic default Ubuntu desktop with kasmVNC
		imageName := "kasmweb/core-kali-rolling:1.18.0-rolling-daily"

		containerID, targetURL, err := helpers.SpawnWorkspaceContainer(imageName, "pwnbox")
		if err != nil {
			http.Error(w, "Failed to spawn workspace: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ws := models.Workspace{
			UserID:      userID,
			ContainerID: containerID,
			Status:      "running",
			TargetURL:   targetURL,
			CreatedAt:   time.Now(),
			ExpiresAt:   time.Now().Add(2 * time.Hour),
		}

		if err := db.Create(&ws).Error; err != nil {
			helpers.StopWorkspaceContainer(containerID)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ws)
	}
}

// StopWorkspace stops and clears the workspace
func StopWorkspace(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var existing models.Workspace
		if err := db.Where("user_id = ?", userID).First(&existing).Error; err != nil {
			http.Error(w, "No active workspace", http.StatusNotFound)
			return
		}

		// Kill the container
		if err := helpers.StopWorkspaceContainer(existing.ContainerID); err != nil {
			// Proceed to delete DB record anyway to avoid locking the user out forever if docker daemon acts up
			helpers.Warning("Failed to stop container %s: %v", existing.ContainerID, err)
		}

		db.Delete(&existing)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}
}

// StatusWorkspace checks if a user has a workspace
func StatusWorkspace(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var existing models.Workspace
		if err := db.Where("user_id = ?", userID).First(&existing).Error; err != nil {
			http.Error(w, "No active workspace", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existing)
	}
}

// WorkspaceProxy acts as a reverse proxy directly into the user's allocated container workspace.
// This naturally supports WebSocket upgrades (noVNC websockify protocol).
func WorkspaceProxy(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var ws models.Workspace
		if err := db.Where("user_id = ?", userID).First(&ws).Error; err != nil {
			http.Error(w, "No active workspace found", http.StatusNotFound)
			return
		}

		target, err := url.Parse(ws.TargetURL)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		// Kasm containers use self-signed certificates for their internal HTTPS/WSS service.
		// We tell the proxy to skip verification for these internal bridge connections.
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		d := proxy.Director
		proxy.Director = func(req *http.Request) {
			d(req)
			req.Host = target.Host // Some VNC / websockify instances require original host matching internally

			// Inject Kasm Basic Auth credentials to bypass the login prompt.
			// Image defaults: user 'kasm_user', password we set 'password'
			auth := "kasm_user:password"
			basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
			req.Header.Set("Authorization", basicAuth)
		}

		proxy.ModifyResponse = func(resp *http.Response) error {
			helpers.Information("[PROXY] Received %d from container", resp.StatusCode)
			return nil
		}

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			helpers.Warning("[PROXY] Error: %v", err)
			http.Error(w, "Proxy Error: "+err.Error(), http.StatusBadGateway)
		}

		proxy.ServeHTTP(w, r)
	}
}

func WorkspaceRoutes(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.AuthMiddleware)

	r.Options("/*", helpers.CORSOptionsHandler)

	r.Post("/start", StartWorkspace(db))
	r.Post("/stop", StopWorkspace(db))
	r.Get("/status", StatusWorkspace(db))

	// The frontend uses /workspace/proxy to load the VNC HTML iframe / Websocket
	r.Handle("/proxy/*", http.StripPrefix("/workspace/proxy", WorkspaceProxy(db)))
	r.Handle("/proxy", http.StripPrefix("/workspace/proxy", WorkspaceProxy(db)))

	return r
}
