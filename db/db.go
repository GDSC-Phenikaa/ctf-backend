package db

import (
	"os"

	"github.com/GDSC-Phenikaa/twilight-ctf/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db_type := os.Getenv("DB_TYPE")

	if db_type == "sqlite" {
		db, err := gorm.Open(sqlite.Open(os.Getenv("DB_NAME")), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		err = db.AutoMigrate(&models.User{})
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	if db_type == "postgres" {
		dsn := os.Getenv("DB_DSN")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		err = db.AutoMigrate(&models.User{})
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return nil, nil
}
