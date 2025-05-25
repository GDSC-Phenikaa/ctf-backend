package models

import "gorm.io/gorm"

type Challanges struct {
	gorm.Model  `swaggerignore:"true"`
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Difficulty  string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Points      int    `gorm:"not null"`
	Flag        string `gorm:"not null"` // The flag for the challenge, which users need to submit to solve it
	CreatedAt   string `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   string `gorm:"not null;default:CURRENT_TIMESTAMP"`
	AuthorID    uint   `gorm:"not null"` // Foreign key to the user who created the challenge
	Author      User   `gorm:"foreignKey:AuthorID;references:ID"`
	AuthorName  string `gorm:"-" json:"author_name"` // Not stored in DB, for response only
	Docker      bool   `gorm:"not null"`             // Indicates if the challenge requires Docker
	DockerImage string
	Solves      int  `gorm:"default:0"`     // Number of times the challenge has been solved
	Hidden      bool `gorm:"default:false"` // Indicates if the challenge is hidden from the public
}
