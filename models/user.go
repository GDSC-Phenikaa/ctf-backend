package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `swaggerignore:"true"`
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Email      string `gorm:"unique;not null"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	IsAdmin    bool   `gorm:"default:false"`
}
