package config

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

// GetDB returns the database connection
func GetDB() *sql.DB {
	once.Do(func() {
		// Ensure `data/` directory exists
		if _, err := os.Stat("./data"); os.IsNotExist(err) {
			os.Mkdir("./data", os.ModePerm)
		}

		var err error
		// ✅ Use a fresh database file
		db, err = sql.Open("sqlite3", "./data/forum_new.db?_busy_timeout=30000&_journal_mode=WAL&_locking_mode=NORMAL")
		if err != nil {
			log.Fatal("Failed to open database:", err)
		}

		// ✅ Enable WAL mode explicitly
		_, err = db.Exec("PRAGMA journal_mode = WAL;")
		if err != nil {
			log.Fatal("Failed to enable WAL mode:", err)
		}

		// ✅ Increase busy timeout
		_, err = db.Exec("PRAGMA busy_timeout = 30000;")
		if err != nil {
			log.Fatal("Failed to set busy timeout:", err)
		}
	})
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed.")
	}
}
\