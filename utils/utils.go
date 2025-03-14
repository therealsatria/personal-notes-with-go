package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error definitions
var (
	ErrCategoryNameEmpty    = errors.New("category name cannot be empty")
	ErrCategoryNameConflict = errors.New("category name already exists")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrNoteSubjectEmpty     = errors.New("note subject cannot be empty")
	ErrNoteNotFound         = errors.New("note not found")
	ErrEmptyInput           = errors.New("input text cannot be empty")
)

// HandleBadRequestError handles bad request errors.
func HandleBadRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

// HandleInternalServerError handles internal server errors.
func HandleInternalServerError(c *gin.Context, err error, operation string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to %s: %v", operation, err)})
}
