package repositories

import (
	"database/sql"
	"fmt"
	"personal-notes-with-go/models"
	"personal-notes-with-go/utils"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type CategoryRepositoryInterface interface {
	Create(category *models.Category) error
	GetAll() ([]models.Category, error)
	GetByID(id string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id string) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepositoryInterface {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	category.ID = uuid.New().String()

	// Encrypt name
	encryptedName, err := utils.Encrypt(category.Name)
	if err != nil {
		return fmt.Errorf("failed to encrypt category name: %w", err)
	}

	_, err = r.db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", category.ID, encryptedName)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return utils.ErrCategoryNameConflict
			}
		}
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		var encryptedName string
		if err := rows.Scan(&cat.ID, &encryptedName); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		// Decrypt name
		cat.Name, err = utils.Decrypt(encryptedName)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt category name: %w", err)
		}

		categories = append(categories, cat)
	}
	return categories, nil
}

func (r *categoryRepository) GetByID(id string) (*models.Category, error) {
	var category models.Category
	var encryptedName string
	err := r.db.QueryRow("SELECT id, name FROM categories WHERE id = ?", id).Scan(&category.ID, &encryptedName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Decrypt name
	category.Name, err = utils.Decrypt(encryptedName)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt category name: %w", err)
	}

	return &category, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	// Encrypt name
	encryptedName, err := utils.Encrypt(category.Name)
	if err != nil {
		return fmt.Errorf("failed to encrypt category name: %w", err)
	}

	result, err := r.db.Exec("UPDATE categories SET name = ? WHERE id = ?", encryptedName, category.ID)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return utils.ErrCategoryNameConflict
			}
		}
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return utils.ErrCategoryNotFound
	}

	return nil
}

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
		return utils.ErrCategoryNotFound
	}

	return nil
}
