package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// GetDB initializes and returns the database connection
func GetDB() *sql.DB {
	if db != nil {
		return db
	}

	// ✅ Ensure `data/` directory exists
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", os.ModePerm)
	}

	var err error
	db, err = sql.Open("sqlite3", "./data/forum.db?_busy_timeout=10000&_journal_mode=WAL&_locking_mode=NORMAL")
	if err != nil {
		log.Fatal("❌ Failed to open database:", err)
	}

	// ✅ Apply migrations before returning DB
	if err := applyMigrations(); err != nil {
		log.Fatal("❌ Failed to apply migrations:", err)
	}

	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// applyMigrations runs all `.up.sql` files inside `migrations/`
func applyMigrations() error {
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
			fmt.Println("🔹 Applied migration:", file.Name())
		}
	}
	return nil
}
