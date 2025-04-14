package db

import (
	"database/sql"
	"fmt"
	"log"
	"ops_service/internal/config"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	var db *sql.DB
	var err error
	
	// Try connecting to the database with retries
	maxRetries := 5
	retryDelay := time.Second * 3
	
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to database (attempt %d/%d)...", i+1, maxRetries)
		
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to database")
				return db, nil
			}
		}
		
		log.Printf("Failed to connect to database: %v. Retrying in %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
} 