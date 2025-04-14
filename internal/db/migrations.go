package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations runs all SQL migrations in the migrations directory
func RunMigrations(db *sql.DB) error {
	migrationDir := "migrations"
	
	// Check if migrations directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		log.Printf("Migrations directory '%s' not found. Looking in parent directory...", migrationDir)
		
		// Try parent directory if current directory doesn't have migrations
		parentMigrationDir := filepath.Join("..", migrationDir)
		if _, err := os.Stat(parentMigrationDir); os.IsNotExist(err) {
			log.Printf("Migrations directory '%s' not found either. Skipping migrations.", parentMigrationDir)
			return nil
		}
		
		migrationDir = parentMigrationDir
	}

	// Get all migration files
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort up migrations
	var upMigrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upMigrations = append(upMigrations, file.Name())
		}
	}
	
	sort.Strings(upMigrations)
	
	if len(upMigrations) == 0 {
		log.Println("No migrations found")
		return nil
	}

	// Run each migration in a transaction
	for _, migrationFile := range upMigrations {
		log.Printf("Running migration: %s", migrationFile)
		
		// Read migration file
		path := filepath.Join(migrationDir, migrationFile)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", migrationFile, err)
		}

		// Execute migration
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migrationFile, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for migration %s: %w", migrationFile, err)
		}
		
		log.Printf("Migration %s completed successfully", migrationFile)
	}

	log.Printf("All migrations completed successfully")
	return nil
} 