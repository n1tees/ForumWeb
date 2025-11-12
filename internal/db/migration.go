package db

import (
	"ForumWeb/internal/db/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.Answer{},
	)
}
