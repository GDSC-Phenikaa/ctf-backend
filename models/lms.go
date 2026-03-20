package models

import "gorm.io/gorm"

type Module struct {
	gorm.Model  `swaggerignore:"true"`
	ID          uint     `gorm:"primaryKey" json:"id"`
	Title       string   `gorm:"not null" json:"title"`
	Description string   `gorm:"not null" json:"description"`
	Order       int      `gorm:"not null" json:"order"`
	Lessons     []Lesson `gorm:"foreignKey:ModuleID" json:"lessons"`
}

type Lesson struct {
	gorm.Model `swaggerignore:"true"`
	ID         uint       `gorm:"primaryKey" json:"id"`
	ModuleID   uint       `gorm:"not null" json:"module_id"`
	Title      string     `gorm:"not null" json:"title"`
	Content    string     `gorm:"not null" json:"content"` // Markdown or HTML
	Order      int        `gorm:"not null" json:"order"`
	Questions  []Question `gorm:"foreignKey:LessonID" json:"questions"`
}

type Question struct {
	gorm.Model    `swaggerignore:"true"`
	ID            uint   `gorm:"primaryKey" json:"id"`
	LessonID      uint   `gorm:"not null" json:"lesson_id"`
	Content       string `gorm:"not null" json:"content"`
	Type          string `gorm:"not null;default:'mcq'" json:"type"` // "mcq" or "text"
	Options       string `gorm:"type:text" json:"options"` // JSON array or comma separated
	CorrectAnswer string `gorm:"not null" json:"correct_answer"`
	Points        int    `gorm:"not null;default:0" json:"points"`
}

type QuestionSolve struct {
	gorm.Model      `swaggerignore:"true"`
	ID              uint     `gorm:"primaryKey"`
	QuestionID      uint     `gorm:"not null"`
	Question        Question `gorm:"foreignKey:QuestionID;references:ID"`
	UserID          uint     `gorm:"not null"`
	User            User     `gorm:"foreignKey:UserID;references:ID"`
	SubmittedAnswer string   `gorm:"not null"`
	Correct         bool     `gorm:"not null"`
}
