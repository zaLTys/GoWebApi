package migrations

import "gorm.io/gorm"

// MigrationManager defines the interface for database migrations
type MigrationManager interface {
	RunMigrations(db *gorm.DB) error
}
