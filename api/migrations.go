package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{db: db}
}

// Run executes all pending migrations
func (m *MigrationRunner) Run(migrationsDir string) error {
	log.Println("[Migrations] Starting migration process...")

	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of already-run migrations
	executedMigrations, err := m.getExecutedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	// Read migration files
	files, err := m.getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		log.Println("[Migrations] No migration files found")
		return nil
	}

	// Execute pending migrations
	executed := 0
	for _, file := range files {
		filename := filepath.Base(file)

		// Skip if already executed
		if executedMigrations[filename] {
			log.Printf("[Migrations] Skipping already executed: %s", filename)
			continue
		}

		log.Printf("[Migrations] Executing: %s", filename)
		if err := m.executeMigration(file, filename); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
		executed++
	}

	if executed == 0 {
		log.Println("[Migrations] No pending migrations")
	} else {
		log.Printf("[Migrations] Successfully executed %d migration(s)", executed)
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *MigrationRunner) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT NOT NULL UNIQUE,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

// getExecutedMigrations returns a map of already-executed migrations
func (m *MigrationRunner) getExecutedMigrations() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT filename FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		executed[filename] = true
	}

	return executed, rows.Err()
}

// getMigrationFiles returns sorted list of migration files
func (m *MigrationRunner) getMigrationFiles(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return nil, err
	}

	// Sort files alphabetically (they should be prefixed with numbers)
	sort.Strings(files)
	return files, nil
}

// executeMigration runs a single migration file
func (m *MigrationRunner) executeMigration(filepath, filename string) error {
	// Read migration file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Split by semicolons to handle multiple statements
	statements := strings.Split(string(content), ";")

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %w\nStatement: %s", err, stmt)
		}
	}

	// Record migration
	_, err = tx.Exec("INSERT INTO migrations (filename) VALUES (?)", filename)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
