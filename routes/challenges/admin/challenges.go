package admin

import (
	"encoding/json"
	"net/http"

	"github.com/GDSC-Phenikaa/ctf-backend/middlewares"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"gorm.io/gorm"
)

func listChallengesHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var challenges []models.Challanges
		if err := database.Preload("Author").Find(&challenges).Error; err != nil {
			http.Error(w, "Failed to retrieve challenges", http.StatusInternalServerError)
			return
		}

		// Set AuthorName for each challenge
		for i := range challenges {
			challenges[i].AuthorName = challenges[i].Author.Name
		}

		for i := range challenges {
			challenges[i].AuthorName = challenges[i].Author.Name
			challenges[i].Author = models.User{} // Do not return Author struct
		}

		response := map[string]interface{}{
			"status":     "success",
			"challenges": challenges,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func createChallengeHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var challenge models.Challanges
		if err := json.NewDecoder(r.Body).Decode(&challenge); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID, ok := middlewares.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		challenge.AuthorID = userID

		if err := database.Create(&challenge).Error; err != nil {
			http.Error(w, "Failed to create challenge", http.StatusInternalServerError)
			return
		}

		var author models.User
		if err := database.First(&author, userID).Error; err == nil {
			challenge.AuthorName = author.Name
		}

		challenge.Author = models.User{} // Do not return Author struct

		response := map[string]interface{}{
			"status":    "success",
			"message":   "Challenge created successfully",
			"challenge": challenge,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func updateChallengeHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var challenge models.Challanges
		if err := json.NewDecoder(r.Body).Decode(&challenge); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := database.Save(&challenge).Error; err != nil {
			http.Error(w, "Failed to update challenge", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"status":    "success",
			"message":   "Challenge updated successfully",
			"challenge": challenge,
		}
		json.NewEncoder(w).Encode(response)
	}
}
