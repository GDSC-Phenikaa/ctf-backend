package models

import (
	"time"
)

type Workspace struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex" json:"user_id"`
	ContainerID string    `gorm:"uniqueIndex;not null" json:"container_id"`
	Status      string    `json:"status"`
	TargetURL   string    `json:"target_url"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
