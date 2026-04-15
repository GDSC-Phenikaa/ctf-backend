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
	gorm.Model  `swaggerignore:"true"`
	ID          uint           `gorm:"primaryKey" json:"id"`
	ModuleID    uint           `gorm:"not null" json:"module_id"`
	Title       string         `gorm:"not null" json:"title"`
	Content     string         `gorm:"not null" json:"content"` // Markdown or HTML
	Body        string         `gorm:"type:text" json:"body"`
	VideoIframe string         `gorm:"type:text" json:"video_iframe"`
	Order       int            `gorm:"not null" json:"order"`
	Questions   []Question     `gorm:"foreignKey:LessonID" json:"questions"`
	Segments    []VideoSegment `gorm:"foreignKey:LessonID" json:"segments"`
}

type VideoSegment struct {
	gorm.Model   `swaggerignore:"true"`
	ID           uint   `gorm:"primaryKey" json:"id"`
	LessonID     uint   `gorm:"not null" json:"lesson_id"`
	Title        string `gorm:"not null" json:"title"`
	Description  string `gorm:"type:text" json:"description"`
	StartSeconds int    `gorm:"not null;default:0" json:"start_seconds"`
	EndSeconds   int    `gorm:"not null;default:0" json:"end_seconds"`
	Order        int    `gorm:"not null;default:0" json:"order"`
}

type Question struct {
	gorm.Model     `swaggerignore:"true"`
	ID             uint   `gorm:"primaryKey" json:"id"`
	LessonID       uint   `gorm:"not null" json:"lesson_id"`
	VideoSegmentID *uint  `json:"video_segment_id"`
	Placement      string `gorm:"not null;default:'lesson'" json:"placement"`
	Content        string `gorm:"not null" json:"content"`
	Prompt         string `gorm:"type:text" json:"prompt"`
	Type           string `gorm:"not null;default:'mcq'" json:"type"` // legacy: "mcq"/"text", v2: single_choice, multi_choice, true_false, short_text, long_text, numeric, code
	Options        string `gorm:"type:text" json:"options"`
	CorrectAnswer  string `gorm:"not null" json:"correct_answer"`
	AnswerKey      string `gorm:"type:text" json:"answer_key"`
	Points         int    `gorm:"not null;default:0" json:"points"`
	Order          int    `gorm:"not null;default:0" json:"order"`
}

type QuestionSolve struct {
	gorm.Model       `swaggerignore:"true"`
	ID               uint     `gorm:"primaryKey"`
	QuestionID       uint     `gorm:"not null"`
	Question         Question `gorm:"foreignKey:QuestionID;references:ID"`
	UserID           uint     `gorm:"not null"`
	User             User     `gorm:"foreignKey:UserID;references:ID"`
	SubmittedAnswer  string   `gorm:"not null"`
	NormalizedAnswer string   `gorm:"type:text" json:"normalized_answer"`
	AttemptNo        int      `gorm:"not null;default:1" json:"attempt_no"`
	Correct          bool     `gorm:"not null"`
}
