package repositories

import (
	"database/sql"
	"fmt"
	"personal-notes-with-go/models"
	"personal-notes-with-go/utils"
	"strings"

	"github.com/google/uuid"
)

// NoteRepository interface mendefinisikan operasi yang dapat dilakukan pada data Note.
type NoteRepository interface {
	Create(note *models.Note) error
	GetAll(priority, categoryID string) ([]models.Note, error)
	Update(note *models.Note) error
	Delete(id string) error
}

// noteRepository struct mengimplementasikan NoteRepository interface.
type noteRepository struct {
	db *sql.DB
}

// NewNoteRepository membuat instance baru dari NoteRepository.
func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{db: db}
}

// Create menambahkan catatan baru ke database.
func (r *noteRepository) Create(note *models.Note) error {
	note.ID = uuid.New().String()
	_, err := r.db.Exec("INSERT INTO notes (id, subject, content, priority, tags, category_id) VALUES (?, ?, ?, ?, ?, ?)",
		note.ID, note.Subject, note.Content, note.Priority, note.Tags, note.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}
	return nil
}

// GetAll mengambil semua catatan dari database, dengan filter opsional.
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
		if err := rows.Scan(&note.ID, &note.Subject, &note.Content, &note.Priority, &note.Tags, &note.CategoryID); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// Update memperbarui catatan yang ada di database.
func (r *noteRepository) Update(note *models.Note) error {
	result, err := r.db.Exec("UPDATE notes SET subject = ?, content = ?, priority = ?, tags = ?, category_id = ? WHERE id = ?",
		note.Subject, note.Content, note.Priority, note.Tags, note.CategoryID, note.ID)
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

// Delete menghapus catatan dari database berdasarkan ID.
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
