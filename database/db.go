package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	// Buat tabel categories
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS categories (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL UNIQUE
    )`)
	if err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}

	// Buat tabel notes
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notes (
        id TEXT PRIMARY KEY,
        subject TEXT NOT NULL,
        content TEXT,
        priority TEXT,
        tags TEXT,
        category_id TEXT,
        FOREIGN KEY (category_id) REFERENCES categories(id)
    )`)
	if err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}
	log.Println("Database and tables created successfully")
	return nil
}
