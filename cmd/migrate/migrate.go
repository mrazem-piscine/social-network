package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// ✅ Ensure database connection is properly opened
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	defer db.Close() // ✅ Ensure database closes when done

	// ✅ Apply Migrations
	err = applyMigrations(db)
	if err != nil {
		log.Fatalf("❌ Failed to apply migrations: %v", err)
	}
	fmt.Println("✅ Migrations applied successfully.")
}

// applyMigrations ensures all tables exist
func applyMigrations(db *sql.DB) error {
	migrationDir := "migrations"
	absPath, err := filepath.Abs(migrationDir)
	if err != nil {
		return fmt.Errorf("❌ Failed to get absolute migration path: %v", err)
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to read migrations directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationPath := filepath.Join(absPath, file.Name())
			migrationSQL, err := os.ReadFile(migrationPath)
			if err != nil {
				return fmt.Errorf("❌ Failed to read migration %s: %v", file.Name(), err)
			}

			fmt.Println("🔹 Executing Migration:", file.Name())
			fmt.Println("🔹 SQL Query:\n", string(migrationSQL)) // ✅ Debugging output

			_, err = db.Exec(string(migrationSQL))
			if err != nil {
				return fmt.Errorf("❌ Failed to execute migration %s: %v", file.Name(), err)
			}
			fmt.Println("✅ Applied migration:", file.Name())
		}
	}
	return nil
}
