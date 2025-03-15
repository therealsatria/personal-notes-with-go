package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

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

	// Buat tabel activity_logs
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS activity_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME NOT NULL,
        action TEXT NOT NULL,
        entity_type TEXT NOT NULL,
        entity_id INTEGER NOT NULL,
        description TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        ip_address TEXT NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_activity_logs_timestamp ON activity_logs(timestamp);
    CREATE INDEX IF NOT EXISTS idx_activity_logs_entity_type ON activity_logs(entity_type);
    CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
    `)
	if err != nil {
		return fmt.Errorf("failed to create activity_logs table: %w", err)
	}

	log.Println("Database and tables created successfully")
	return nil
}

// FixEncryptionIssues attempts to fix any encryption issues in the database
// This function should be called after encryption is initialized
func FixEncryptionIssues(db *sql.DB) error {
	log.Println("Checking for encryption issues in the database...")

	// Check if there are any notes in the database
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count notes: %w", err)
	}

	if count == 0 {
		log.Println("No notes found in the database, skipping encryption check")
		return nil
	}

	// Get all notes
	rows, err := db.Query("SELECT id, subject, content, tags FROM notes")
	if err != nil {
		return fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	// Check each note for encryption issues
	for rows.Next() {
		var id, subject, content, tags string
		err := rows.Scan(&id, &subject, &content, &tags)
		if err != nil {
			return fmt.Errorf("failed to scan note: %w", err)
		}

		// Check if any of the fields are not base64 encoded
		needsUpdate := false
		if !isBase64(subject) || !isBase64(content) || !isBase64(tags) {
			log.Printf("Found note with ID %s that has encryption issues", id)
			needsUpdate = true
		}

		if needsUpdate {
			// Delete the note
			_, err := db.Exec("DELETE FROM notes WHERE id = ?", id)
			if err != nil {
				return fmt.Errorf("failed to delete note with encryption issues: %w", err)
			}
			log.Printf("Deleted note with ID %s due to encryption issues", id)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating notes: %w", err)
	}

	log.Println("Encryption issues check completed")
	return nil
}

// isBase64 checks if a string is base64 encoded
func isBase64(s string) bool {
	// Empty strings are considered valid
	if s == "" {
		return true
	}

	// Check if the string is base64 encoded
	s = strings.TrimSpace(s)
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
