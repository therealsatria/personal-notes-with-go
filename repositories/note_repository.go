package repositories
import (
	"os"
)

import (
	"database/sql"
	"fmt"
	"personal-notes-with-go/models"
	"personal-notes-with-go/utils"
	"strings"

	"github.com/google/uuid"
)

type NoteRepository interface {
	Create(note *models.Note) error
	GetAll(priority, categoryID string) ([]models.Note, error)
	Update(note *models.Note) error
	Delete(id string) error
}

type noteRepository struct {
	db          *sql.DB
	keyFilePath string // Tambahkan keyFilePath
}

func NewNoteRepository(db *sql.DB) NoteRepository {
	// Pastikan file kunci enkripsi ada
	keyFilePath := "./encryption.key"
	if _, err := os.Stat(keyFilePath); err != nil {
		panic(fmt.Sprintf("Encryption key file not found: %v", err))
	}
	return &noteRepository{db: db, keyFilePath: keyFilePath}
}

func (r *noteRepository) Create(note *models.Note) error {
	note.ID = uuid.New().String()

	// Enkripsi Subject dan Content
	encryptedSubject, err := utils.EncryptString(note.Subject, r.keyFilePath)
	if err != nil {
		return fmt.Errorf("failed to encrypt subject: %w", err)
	}
	encryptedContent, err := utils.EncryptString(note.Content, r.keyFilePath)
	if err != nil {
		return fmt.Errorf("failed to encrypt content: %w", err)
	}

	_, err = r.db.Exec("INSERT INTO notes (id, subject, content, priority, tags, category_id) VALUES (?, ?, ?, ?, ?, ?)",
		note.ID, encryptedSubject, encryptedContent, note.Priority, note.Tags, note.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}
	return nil
}

func (r *noteRepository) GetAll(priority, categoryID string) ([]models.Note, error) {
	query := "SELECT id, subject, content, priority, tags, category_id FROM notes"
	var args []interface{}
	var conditions []string

	if priority != "" {
		conditions = append(conditions, "priority = ?")
		args = append(args, priority)
	}
	if categoryID != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, categoryID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		var encryptedSubject, encryptedContent string
		if err := rows.Scan(&note.ID, &encryptedSubject, &encryptedContent, &note.Priority, &note.Tags, &note.CategoryID); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		// Dekripsi Subject dan Content
		note.Subject, err = utils.DecryptString(encryptedSubject, r.keyFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt subject: %w", err)
		}
		note.Content, err = utils.DecryptString(encryptedContent, r.keyFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt content: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (r *noteRepository) Update(note *models.Note) error {
	// Enkripsi Subject dan Content
	encryptedSubject, err := utils.EncryptString(note.Subject, r.keyFilePath)
	if err != nil {
		return fmt.Errorf("failed to encrypt subject: %w", err)
	}
	encryptedContent, err := utils.EncryptString(note.Content, r.keyFilePath)
	if err != nil {
		return fmt.Errorf("failed to encrypt content: %w", err)
	}

	result, err := r.db.Exec("UPDATE notes SET subject = ?, content = ?, priority = ?, tags = ?, category_id = ? WHERE id = ?",
		encryptedSubject, encryptedContent, note.Priority, note.Tags, note.CategoryID, note.ID)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return utils.ErrNoteNotFound
	}

	return nil
}

func (r *noteRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return utils.ErrNoteNotFound
	}

	return nil
}
