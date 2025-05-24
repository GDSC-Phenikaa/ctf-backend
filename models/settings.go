package models

import "gorm.io/gorm"

type Settings struct {
	gorm.Model  `swaggerignore:"true"`
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	StartTime   string `gorm:"not null"`
	EndTime     string `gorm:"not null"`
}
