package models

import "gorm.io/gorm"

type Container struct {
	gorm.Model `swaggerignore:"true"`
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Image      string `gorm:"not null"`
	Port       int    `gorm:"not null"`
	Env        string `gorm:"not null"`
	Network    string `gorm:"not null"`
	Volume     string `gorm:"not null"`
	Challenge  Challanges `gorm:"not null"`
	ContainerID string `gorm:"not null"`
}