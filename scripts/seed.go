package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GDSC-Phenikaa/ctf-backend/db"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	fmt.Println("Generating more fake data...")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	// Admin
	var admin models.User
	if err := database.Where("username = ?", "admin").First(&admin).Error; err != nil {
		admin = models.User{Name: "Admin", Username: "admin", Email: "admin@ctf.com", Password: string(passwordHash), IsAdmin: true}
		database.Create(&admin)
	}

	// Users
	users := []models.User{
		{Name: "Alice Hacker", Username: "alice", Email: "alice@ctf.com", Password: string(passwordHash), IsAdmin: false},
		{Name: "Bob Exploit", Username: "bob", Email: "bob@ctf.com", Password: string(passwordHash), IsAdmin: false},
		{Name: "Charlie Root", Username: "charlie", Email: "charlie@ctf.com", Password: string(passwordHash), IsAdmin: false},
		{Name: "Diana Pwn", Username: "diana", Email: "diana@ctf.com", Password: string(passwordHash), IsAdmin: false},
		{Name: "Eve Sniff", Username: "eve", Email: "eve@ctf.com", Password: string(passwordHash), IsAdmin: false},
		{Name: "Frank Buffer", Username: "frank", Email: "frank@ctf.com", Password: string(passwordHash), IsAdmin: false},
	}

	for i, u := range users {
		var existing models.User
		if err := database.Where("username = ?", u.Username).First(&existing).Error; err != nil {
			database.Create(&u)
			users[i] = u
		} else {
			users[i] = existing
		}
	}

	// Map of usernames to IDs
	userIDs := make(map[string]uint)
	for _, u := range users {
		userIDs[u.Username] = u.ID
	}

	// CTF Challenges
	challenges := []models.Challanges{
		{Title: "Baby Web", Description: "Find the flag in the source code.", Difficulty: "Easy", Type: "Web", Points: 100, Flag: "GDSC{baby_web_123}", AuthorID: admin.ID, Docker: false},
		{Title: "SQL Injection 101", Description: "Bypass the login.", Difficulty: "Medium", Type: "Web", Points: 250, Flag: "GDSC{sqli_is_fun}", AuthorID: admin.ID, Docker: false},
		{Title: "RSA Basics", Description: "Given p, q, and e, find d.", Difficulty: "Medium", Type: "Crypto", Points: 300, Flag: "GDSC{math_magic}", AuthorID: admin.ID, Docker: false},
		{Title: "Buffer Overflow 1", Description: "Smash the stack.", Difficulty: "Hard", Type: "Pwn", Points: 500, Flag: "GDSC{ret2win}", AuthorID: admin.ID, Docker: true},
	}

	for i, c := range challenges {
		var existing models.Challanges
		if err := database.Where("title = ?", c.Title).First(&existing).Error; err != nil {
			database.Create(&c)
			challenges[i] = c
		} else {
			challenges[i] = existing
		}
	}

	// Add CTF Solves
	solveMap := map[string][]string{
		"Baby Web":          {"alice", "bob", "diana", "eve", "frank"}, // 100
		"SQL Injection 101": {"alice", "eve", "frank"},                 // 250
		"RSA Basics":        {"bob", "diana"},                          // 300
		"Buffer Overflow 1": {"alice", "diana"},                        // 500
	}

	for _, c := range challenges {
		for _, username := range solveMap[c.Title] {
			createCTFSolve(database, c.ID, userIDs[username], true, c.Flag)
		}
	}

	// LMS Modules, Lessons, Questions
	modules := []models.Module{
		{Title: "Intro to Security", Description: "Basics of cybersecurity", Order: 1},
		{Title: "Web Exploitation", Description: "Finding and exposing logic flaws in web applications", Order: 2},
		{Title: "Reverse Engineering", Description: "Decompiling and dissecting software", Order: 3},
		{Title: "Network Security", Description: "Understanding and securing computer networks", Order: 4},
		{Title: "Forensics", Description: "Digital forensics and incident response", Order: 5},
	}

	for i, m := range modules {
		var existing models.Module
		if err := database.Where("title = ?", m.Title).First(&existing).Error; err != nil {
			database.Create(&m)
			modules[i] = m
		} else {
			modules[i] = existing
		}
	}

	lessons := []models.Lesson{
		// Module 1 (Intro to Security)
		{ModuleID: modules[0].ID, Title: "What is a CTF?", Content: "# Capture the Flag\nCTF is a cybersecurity competition...", Order: 1},
		{ModuleID: modules[0].ID, Title: "Cryptography basics", Content: "# Encryption\nEncryption is the process of encoding information...", Order: 2},
		
		// Module 2 (Web Exploitation)
		{ModuleID: modules[1].ID, Title: "Introduction to HTML/JS", Content: "Web browsers run HTML for layout and JS for logic...", Order: 1},
		{ModuleID: modules[1].ID, Title: "SQL Injection", Content: "SQLi happens when user input is incorrectly parsed as a SQL command...", Order: 2},
		{ModuleID: modules[1].ID, Title: "Cross-Site Scripting (XSS)", Content: "XSS vulnerabilities allow an attacker to inject dangerous Javascript...", Order: 3},
		
		// Module 3 (Reverse Engineering)
		{ModuleID: modules[2].ID, Title: "Assembly 101", Content: "Computers execute machine code, which is assembled from Assembly instructions like MOV, ADD, JMP.", Order: 1},
		{ModuleID: modules[2].ID, Title: "GDB and Debuggers", Content: "Debuggers allow you to pause execution and inspect memory contents.", Order: 2},

		// Module 4 (Network Security)
		{ModuleID: modules[3].ID, Title: "OSI Model", Content: "The OSI model conceptually divides networks into 7 layers...", Order: 1},
		{ModuleID: modules[3].ID, Title: "Packet Sniffing with Wireshark", Content: "Wireshark allows you to capture raw packets traveling across your network.", Order: 2},
		
		// Module 5 (Forensics)
		{ModuleID: modules[4].ID, Title: "File Signatures", Content: "Magic bytes indicate the true file format regardless of the extension.", Order: 1},
		{ModuleID: modules[4].ID, Title: "Steganography", Content: "Hiding data within other data, typically images or audio files.", Order: 2},
	}

	for i, l := range lessons {
		var existing models.Lesson
		if err := database.Where("title = ?", l.Title).First(&existing).Error; err != nil {
			database.Create(&l)
			lessons[i] = l
		} else {
			lessons[i] = existing
		}
	}

	questions := []models.Question{
		{LessonID: lessons[0].ID, Content: "What is the goal of a CTF?", Type: "mcq", Options: `["To capture the flag", "To write code", "To design graphics"]`, CorrectAnswer: "To capture the flag", Points: 50},
		{LessonID: lessons[0].ID, Content: "What does 'Pwn' mean in CTF terminology?", Type: "text", Options: `[]`, CorrectAnswer: "exploit", Points: 100},
		{LessonID: lessons[1].ID, Content: "Which encryption type uses two different keys?", Type: "mcq", Options: `["Symmetric", "Asymmetric", "Hashing"]`, CorrectAnswer: "Asymmetric", Points: 150},
		
		{LessonID: lessons[2].ID, Content: "Which tag is used for the largest heading in HTML?", Type: "mcq", Options: `["<heading>", "<h6>", "<h1>", "<header>"]`, CorrectAnswer: "<h1>", Points: 50},
		{LessonID: lessons[3].ID, Content: "What does SQL stand for?", Type: "mcq", Options: `["Structured Query Language", "Strong Question Language", "System Query Logic"]`, CorrectAnswer: "Structured Query Language", Points: 75},
		{LessonID: lessons[4].ID, Content: "Which type of XSS is stored on the server's database?", Type: "mcq", Options: `["Reflected XSS", "Stored XSS", "DOM-based XSS"]`, CorrectAnswer: "Stored XSS", Points: 120},
		
		{LessonID: lessons[5].ID, Content: "Which assembly instruction is used to copy data?", Type: "text", Options: `[]`, CorrectAnswer: "MOV", Points: 100},
		{LessonID: lessons[6].ID, Content: "What does GDB stand for?", Type: "text", Options: `[]`, CorrectAnswer: "GNU Debugger", Points: 150},

		{LessonID: lessons[7].ID, Content: "At which OSI layer do IP addresses operate?", Type: "mcq", Options: `["Layer 2 (Data Link)", "Layer 3 (Network)", "Layer 4 (Transport)"]`, CorrectAnswer: "Layer 3 (Network)", Points: 150},
		{LessonID: lessons[8].ID, Content: "What port does HTTP commonly use?", Type: "mcq", Options: `["21", "22", "80", "443"]`, CorrectAnswer: "80", Points: 50},

		{LessonID: lessons[9].ID, Content: "What are the first two bytes of a standard DOS/Windows executable (.exe)?", Type: "text", Options: `[]`, CorrectAnswer: "MZ", Points: 200},
		{LessonID: lessons[10].ID, Content: "What command line tool is commonly used to extract hidden files from images?", Type: "mcq", Options: `["grep", "steghide", "nmap"]`, CorrectAnswer: "steghide", Points: 100},
	}

	for i, q := range questions {
		var existing models.Question
		if err := database.Where("content = ?", q.Content).First(&existing).Error; err != nil {
			database.Create(&q)
			questions[i] = q
		} else {
			questions[i] = existing
		}
	}

	// Add LMS Solves
	lmsSolveMap := map[string][]string{
		"What is the goal of a CTF?":                            {"bob", "charlie", "diana", "eve", "frank"}, 
		"What does 'Pwn' mean in CTF terminology?":              {"charlie", "eve", "frank"},                
		"Which encryption type uses two different keys?":         {"eve", "frank"},                          
		"Which tag is used for the largest heading in HTML?":    {"alice", "bob", "frank"},
		"What does SQL stand for?":                               {"alice", "bob", "charlie", "diana"},
		"Which type of XSS is stored on the server's database?":  {"alice", "diana"},
		"Which assembly instruction is used to copy data?":       {"diana"},
		"What does GDB stand for?":                               {"charlie", "eve"},
		"At which OSI layer do IP addresses operate?":            {"alice", "eve", "frank"},
		"What port does HTTP commonly use?":                      {"alice", "bob", "charlie", "diana", "eve", "frank"},
		"What are the first two bytes of a standard DOS/Windows executable (.exe)?": {"diana", "frank"},
		"What command line tool is commonly used to extract hidden files from images?": {"bob", "eve", "frank"},
	}

	for _, q := range questions {
		for _, username := range lmsSolveMap[q.Content] {
			createLMSSolve(database, q.ID, userIDs[username], true, "dummy answer")
		}
	}

	fmt.Println("Added a ton of new fake data!")
}

func createCTFSolve(database *gorm.DB, challengeID, userID uint, correct bool, flag string) {
	var solve models.Solves
	if err := database.Where("challenge_id = ? AND user_id = ?", challengeID, userID).First(&solve).Error; err != nil {
		database.Create(&models.Solves{ChallengeID: challengeID, UserID: userID, Correct: correct, Flag: flag})
	}
}

func createLMSSolve(database *gorm.DB, questionID, userID uint, correct bool, answer string) {
	var solve models.QuestionSolve
	if err := database.Where("question_id = ? AND user_id = ?", questionID, userID).First(&solve).Error; err != nil {
		database.Create(&models.QuestionSolve{QuestionID: questionID, UserID: userID, Correct: correct, SubmittedAnswer: answer})
	}
}
