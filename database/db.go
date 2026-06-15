package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/urlshortner?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables()
	log.Println("Database connected")
}

func createTables() {
	queries := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		is_verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		original_url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		custom_alias TEXT,
		click_count INTEGER DEFAULT 0,
		last_accessed TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS verification_codes (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		code TEXT NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code);
	CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id);
	`

	if _, err := DB.Exec(queries); err != nil {
		log.Fatal(err)
	}

	DB.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS username TEXT UNIQUE")
	DB.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE")
}
