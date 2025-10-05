package migrations

import (
	"books-api/app/models"
	"log"
	"gorm.io/gorm"
)

// migrationManager implements the MigrationManager interface
type migrationManager struct{}

// NewMigrationManager creates a new instance of migration manager
func NewMigrationManager() MigrationManager {
	return &migrationManager{}
}

// RunMigrations executes all database migrations
func (m *migrationManager) RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migrations...")
	
	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.Book{},
		// Add more models here as your application grows
	)
	
	if err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return err
	}
	
	log.Println("Database migrations completed successfully")
	return nil
}
