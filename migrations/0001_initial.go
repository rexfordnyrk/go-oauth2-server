package migrations

import (
	"fmt"
	"log"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
	"github.com/RichardKnop/go-microservice-example/service"
)

func migrate0001() error {
	migrationName := "0001_initial"

	// Config factory
	cnf := config.NewConfig()

	// Database connection factory
	db, err := database.NewDatabase(cnf)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %s", err)
	}

	migration := &Migration{}
	if err := db.Where(&Migration{Name: migrationName}).First(migration).Error; err != nil {
		log.Printf("Running %s migration", migrationName)

		// Create users table
		if err := db.CreateTable(&service.User{}).Error; err != nil {
			return fmt.Errorf("Error creating users table: %s", db.Error)
		}

		// Create refresh_tokens table
		if err := db.CreateTable(&service.RefreshToken{}).Error; err != nil {
			return fmt.Errorf("Error creating refresh_tokens table: %s", db.Error)
		}

		// Create access_tokens table
		if err := db.CreateTable(&service.AccessToken{}).Error; err != nil {
			return fmt.Errorf("Error creating access_tokens table: %s", db.Error)
		}

		// Add foreign key on access_tokens.user_id
		if err := db.Model(&service.AccessToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on access_tokens.user_id for users(id): %s", db.Error)
		}

		// Add foreign key on access_tokens.refresh_token_id
		if err := db.Model(&service.AccessToken{}).AddForeignKey("refresh_token_id", "refresh_tokens(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on access_tokens.refresh_token_id for refresh_tokens(id): %s", db.Error)
		}

		// Save a record to migrations table,
		// so we don't rerun this migration again
		migration.Name = migrationName
		if err := db.Create(migration).Error; err != nil {
			return fmt.Errorf("Error saving record to migrations table: %s", err)
		}
	} else {
		log.Printf("Skipping %s migration", migrationName)
	}

	return nil
}