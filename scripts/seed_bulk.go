package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/GDSC-Phenikaa/ctf-backend/db"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	fmt.Println("Connecting to DB...")
	database, err := db.Connect()
	if err != nil {
		panic(err)
	}

	fmt.Println("Generating 100+ bulk fake questions...")

	// Fetch all lessons to attach questions to
	var lessons []models.Lesson
	database.Find(&lessons)
	if len(lessons) == 0 {
		log.Fatal("No lessons found to attach questions to!")
	}

	// Fetch all users to attach solves to
	var users []models.User
	database.Where("is_admin = ?", false).Find(&users)
	if len(users) == 0 {
		log.Fatal("No users found to attach solves to!")
	}

	rand.Seed(time.Now().UnixNano())

	// Generate 120 questions
	for i := 1; i <= 120; i++ {
		randomLesson := lessons[rand.Intn(len(lessons))]
		
		qTitle := fmt.Sprintf("Bulk Auto-Generated Question #%d", i)
		qType := "mcq"
		if rand.Float32() > 0.7 {
			qType = "text"
		}
		
		points := (rand.Intn(10) + 1) * 10 // 10 to 100 points

		q := models.Question{
			LessonID:      randomLesson.ID,
			Content:       qTitle,
			Type:          qType,
			Options:       `["Option A", "Option B", "Option C", "Option D"]`,
			CorrectAnswer: "Option A",
			Points:        points,
		}

		// Don't recreate if it exists somehow
		var existing models.Question
		if err := database.Where("content = ?", q.Content).First(&existing).Error; err != nil {
			database.Create(&q)
		} else {
			q = existing
		}

		// Generate random solves around this question
		numSolvers := rand.Intn(len(users) + 1) // 0 to all users
		
		// Shuffle users
		rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
		
		for j := 0; j < numSolvers; j++ {
			createLMSSolve(database, q.ID, users[j].ID, true, "Option A")
		}
	}

	fmt.Println("Successfully seeded 120 bulk questions and random solves!")
}

func createLMSSolve(database *gorm.DB, questionID, userID uint, correct bool, answer string) {
	var solve models.QuestionSolve
	if err := database.Where("question_id = ? AND user_id = ?", questionID, userID).First(&solve).Error; err != nil {
		database.Create(&models.QuestionSolve{QuestionID: questionID, UserID: userID, Correct: correct, SubmittedAnswer: answer})
	}
}
