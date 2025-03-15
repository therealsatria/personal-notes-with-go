package handlers

import (
	"net/http"
	"personal-notes-with-go/models"
	"personal-notes-with-go/repositories"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo           repositories.CategoryRepositoryInterface
	activityLogger *ActivityLogHandler
}

func NewCategoryHandler(repo repositories.CategoryRepositoryInterface) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// SetActivityLogger sets the activity logger for this handler
func (h *CategoryHandler) SetActivityLogger(logger *ActivityLogHandler) {
	h.activityLogger = logger
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.repo.Create(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	// Log activity
	if h.activityLogger != nil {
		description := "Created category: " + category.Name
		h.activityLogger.LogActivity(c, "create", "category", category.ID, description)
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategories returns all categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	// Log activity
	if h.activityLogger != nil {
		description := "Retrieved all categories"
		h.activityLogger.LogActivity(c, "read", "category", 0, description)
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory updates a category by ID
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if category exists
	_, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.ID = id
	if err := h.repo.Update(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	// Log activity
	if h.activityLogger != nil {
		description := "Updated category: " + category.Name
		h.activityLogger.LogActivity(c, "update", "category", id, description)
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category by ID
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Check if category exists
	category, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	// Log activity
	if h.activityLogger != nil {
		description := "Deleted category: " + category.Name
		h.activityLogger.LogActivity(c, "delete", "category", id, description)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
