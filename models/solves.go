package models

type Solves struct {
	ID          uint       `gorm:"primaryKey"`
	ChallengeID uint       `gorm:"not null"` // Foreign key to the challenge that was solved
	Challenge   Challanges `gorm:"foreignKey:ChallengeID;references:ID"`
	UserID      uint       `gorm:"not null"` // Foreign key to the user who solved the challenge
	User        User       `gorm:"foreignKey:UserID;references:ID"`
	Flag        string     `gorm:"not null"` // The flag submitted by the user to solve the challenge
	Correct     bool       `gorm:"not null"` // Indicates if the submitted flag was correct
}
