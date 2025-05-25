package db

import (
	"github.com/GDSC-Phenikaa/ctf-backend/env"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db_type := env.DbType()

	if db_type == "sqlite" {
		db, err := gorm.Open(sqlite.Open(env.DbName()), &gorm.Config{})
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
		dsn := env.DbDsn()
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
