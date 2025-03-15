package repositories

import (
	"database/sql"
	"fmt"
	"personal-notes-with-go/models"
	"personal-notes-with-go/utils"

	"github.com/google/uuid"
)

type NoteRepositoryInterface interface {
	Create(note *models.Note) error
	GetAll() ([]*models.Note, error)
	GetByID(id string) (*models.Note, error)
	Update(note *models.Note) error
	Delete(id string) error
	GetByCategoryID(categoryID string) ([]*models.Note, error)
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) NoteRepositoryInterface {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(note *models.Note) error {
	// Generate a new UUID for the note
	note.ID = uuid.New().String()

	// Note: All encryption is now done in the handler
	// We just insert the already encrypted data

	// Insert into database
	query := `
		INSERT INTO notes (id, subject, content, priority, tags, category_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, note.ID, note.Subject, note.Content, note.Priority, note.Tags, note.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}
	return nil
}

func (r *noteRepository) GetAll() ([]*models.Note, error) {
	query := `SELECT id, subject, content, priority, tags, category_id FROM notes`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		var encryptedSubject, encryptedContent, encryptedTags string
		err := rows.Scan(&note.ID, &encryptedSubject, &encryptedContent, &note.Priority, &encryptedTags, &note.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		// Store encrypted values for handler to decrypt
		note.Subject = encryptedSubject
		note.Content = encryptedContent
		note.Tags = encryptedTags

		notes = append(notes, note)
	}

	return notes, nil
}

func (r *noteRepository) GetByID(id string) (*models.Note, error) {
	query := `SELECT id, subject, content, priority, tags, category_id FROM notes WHERE id = ?`
	note := &models.Note{}
	var encryptedSubject, encryptedContent, encryptedTags string
	err := r.db.QueryRow(query, id).Scan(&note.ID, &encryptedSubject, &encryptedContent, &note.Priority, &encryptedTags, &note.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNoteNotFound
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	// Store encrypted values for handler to decrypt
	note.Subject = encryptedSubject
	note.Content = encryptedContent
	note.Tags = encryptedTags

	return note, nil
}

func (r *noteRepository) Update(note *models.Note) error {
	// Note: All encryption is now done in the handler
	// We just update with the already encrypted data

	query := `
		UPDATE notes
		SET subject = ?, content = ?, priority = ?, tags = ?, category_id = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, note.Subject, note.Content, note.Priority, note.Tags, note.CategoryID, note.ID)
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

// GetByCategoryID returns all notes for a specific category
func (r *noteRepository) GetByCategoryID(categoryID string) ([]*models.Note, error) {
	query := `SELECT id, subject, content, priority, tags, category_id FROM notes WHERE category_id = ?`

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes by category ID: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		note := &models.Note{}
		var encryptedSubject, encryptedContent, encryptedTags string
		err := rows.Scan(&note.ID, &encryptedSubject, &encryptedContent, &note.Priority, &encryptedTags, &note.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note row: %w", err)
		}

		// Store encrypted values for handler to decrypt
		note.Subject = encryptedSubject
		note.Content = encryptedContent
		note.Tags = encryptedTags

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating note rows: %w", err)
	}

	return notes, nil
}
