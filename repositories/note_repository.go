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

	// Encrypt sensitive data
	encryptedContent, err := utils.Encrypt(note.Content)
	if err != nil {
		return fmt.Errorf("failed to encrypt note content: %w", err)
	}

	encryptedTags, err := utils.Encrypt(note.Tags)
	if err != nil {
		return fmt.Errorf("failed to encrypt note tags: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO notes (id, subject, content, priority, tags, category_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.Exec(query, note.ID, note.Subject, encryptedContent, note.Priority, encryptedTags, note.CategoryID)
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
		var encryptedContent, encryptedTags string
		err := rows.Scan(&note.ID, &note.Subject, &encryptedContent, &note.Priority, &encryptedTags, &note.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		// Decrypt sensitive data
		note.Content, err = utils.Decrypt(encryptedContent)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt note content: %w", err)
		}

		note.Tags, err = utils.Decrypt(encryptedTags)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt note tags: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (r *noteRepository) GetByID(id string) (*models.Note, error) {
	query := `SELECT id, subject, content, priority, tags, category_id FROM notes WHERE id = ?`
	note := &models.Note{}
	var encryptedContent, encryptedTags string
	err := r.db.QueryRow(query, id).Scan(&note.ID, &note.Subject, &encryptedContent, &note.Priority, &encryptedTags, &note.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNoteNotFound
		}
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	// Decrypt sensitive data
	note.Content, err = utils.Decrypt(encryptedContent)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt note content: %w", err)
	}

	note.Tags, err = utils.Decrypt(encryptedTags)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt note tags: %w", err)
	}

	return note, nil
}

func (r *noteRepository) Update(note *models.Note) error {
	// Encrypt sensitive data
	encryptedContent, err := utils.Encrypt(note.Content)
	if err != nil {
		return fmt.Errorf("failed to encrypt note content: %w", err)
	}

	encryptedTags, err := utils.Encrypt(note.Tags)
	if err != nil {
		return fmt.Errorf("failed to encrypt note tags: %w", err)
	}

	query := `
		UPDATE notes
		SET subject = ?, content = ?, priority = ?, tags = ?, category_id = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, note.Subject, encryptedContent, note.Priority, encryptedTags, note.CategoryID, note.ID)
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
