package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./urlshortner.db"
	}
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables()
	log.Println("Database initialized")
}

func createTables() {
	queries := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		original_url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		custom_alias TEXT,
		click_count INTEGER DEFAULT 0,
		last_accessed DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code);
	CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id);
	`

	if _, err := DB.Exec(queries); err != nil {
		log.Fatal(err)
	}
}
