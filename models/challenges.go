package models

import "gorm.io/gorm"

type Challanges struct {
	gorm.Model `swaggerignore:"true"`
	ID         uint   `gorm:"primaryKey"`
	Title	   string `gorm:"not null"`
	Description string `gorm:"not null"`
	Flag	   string `gorm:"not null"`
	Points	   int    `gorm:"not null"`
	Category   string `gorm:"not null"`
	Author	   User `gorm:"not null"`
	IsSpawnable bool `gorm:"not null"`
}	
