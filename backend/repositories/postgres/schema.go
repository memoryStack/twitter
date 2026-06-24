package postgres

import (
	"fmt"

	"gorm.io/gorm"

	"twitter/backend/models"
)

// Migrate creates tables and adapter-level constraints (FK, indexes, etc.).
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}, &models.Tweet{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return ensureTweetUserForeignKey(db)
}

func ensureTweetUserForeignKey(db *gorm.DB) error {
	const constraintName = "fk_tweets_user_id"

	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.table_constraints
			WHERE constraint_name = ?
			  AND table_name = 'tweets'
		)
	`, constraintName).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("check tweet user fk: %w", err)
	}
	if exists {
		return nil
	}

	return db.Exec(`
		ALTER TABLE tweets
		ADD CONSTRAINT fk_tweets_user_id
		FOREIGN KEY (user_id) REFERENCES users(id)
		ON UPDATE CASCADE ON DELETE CASCADE
	`).Error
}
