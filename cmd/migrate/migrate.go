package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"social-network/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := config.GetDB()
	defer config.CloseDB()

	// ✅ Ensure migrations are applied separately
	err := applyMigrations(db)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	fmt.Println("Migrations applied successfully.")
}

// applyMigrations ensures all tables exist
func applyMigrations(db *sql.DB) error {
	migrationDir := "migrations"
	absPath, err := filepath.Abs(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute migration path: %v", err)
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationPath := filepath.Join(absPath, file.Name())
			migrationSQL, err := os.ReadFile(migrationPath)
			if err != nil {
				return fmt.Errorf("failed to read migration %s: %v", file.Name(), err)
			}

			_, err = db.Exec(string(migrationSQL))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %v", file.Name(), err)
			}
			fmt.Println("Applied migration:", file.Name())
		}
	}
	return nil
}
