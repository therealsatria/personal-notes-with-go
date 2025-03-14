package repositories

import (
	"database/sql"
	"fmt"
	"personal-notes-with-go/models"
	"personal-notes-with-go/utils"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

// CategoryRepository interface mendefinisikan operasi yang dapat dilakukan pada data Category.
type CategoryRepository interface {
	Create(category *models.Category) error
	GetAll() ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id string) error
}

// categoryRepository struct mengimplementasikan CategoryRepository interface.
type categoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository membuat instance baru dari CategoryRepository.
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create menambahkan kategori baru ke database.
func (r *categoryRepository) Create(category *models.Category) error {
	category.ID = uuid.New().String()
	_, err := r.db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", category.ID, category.Name)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return utils.ErrCategoryNameConflict // Mengembalikan error yang lebih spesifik
			}
		}
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// GetAll mengambil semua kategori dari database.
func (r *categoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// Update memperbarui kategori yang ada di database.
func (r *categoryRepository) Update(category *models.Category) error {
	result, err := r.db.Exec("UPDATE categories SET name = ? WHERE id = ?", category.Name, category.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return utils.ErrCategoryNotFound // Mengembalikan error not found
	}

	return nil
}

// Delete menghapus kategori dari database berdasarkan ID.
func (r *categoryRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return utils.ErrCategoryNotFound // Mengembalikan error not found
	}

	return nil
}
