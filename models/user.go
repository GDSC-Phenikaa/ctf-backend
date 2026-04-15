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

type Note struct {
	gorm.Model `swaggerignore:"true"`
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `gorm:"not null;index" json:"user_id"`
	User       User   `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Href       string `gorm:"not null;index" json:"href"`
	Content    string `gorm:"type:text;not null" json:"content"`
}
