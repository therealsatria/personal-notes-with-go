package handlers

import (
	"fmt"
	"net/http"
	"personal-notes-with-go/models"
	"personal-notes-with-go/repositories"
	"personal-notes-with-go/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	repo           repositories.NoteRepositoryInterface
	activityLogger *ActivityLogHandler
}

func NewNoteHandler(repo repositories.NoteRepositoryInterface) *NoteHandler {
	return &NoteHandler{repo: repo}
}

// SetActivityLogger sets the activity logger for this handler
func (h *NoteHandler) SetActivityLogger(logger *ActivityLogHandler) {
	h.activityLogger = logger
}

// CreateNote creates a new note
func (h *NoteHandler) CreateNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Encrypt sensitive data
	encryptedSubject, err := utils.Encrypt(note.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt subject"})
		return
	}
	note.Subject = encryptedSubject

	encryptedContent, err := utils.Encrypt(note.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt content"})
		return
	}
	note.Content = encryptedContent

	encryptedTags, err := utils.Encrypt(note.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt tags"})
		return
	}
	note.Tags = encryptedTags

	if err := h.repo.Create(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	// Log the activity
	if h.activityLogger != nil {
		noteID, _ := strconv.Atoi(note.ID)
		decryptedSubject, _ := utils.Decrypt(note.Subject)
		h.activityLogger.LogActivity(c, "create", "note", noteID, "Created note: "+decryptedSubject)
	}

	// Decrypt for response
	decryptedSubject, err := utils.Decrypt(note.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt subject"})
		return
	}
	note.Subject = decryptedSubject

	decryptedContent, err := utils.Decrypt(note.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt content"})
		return
	}
	note.Content = decryptedContent

	decryptedTags, err := utils.Decrypt(note.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt tags"})
		return
	}
	note.Tags = decryptedTags

	c.JSON(http.StatusCreated, note)
}

// GetNotes returns all notes or filtered by category
func (h *NoteHandler) GetNotes(c *gin.Context) {
	categoryID := c.Query("category_id")
	limit := 10 // Default limit

	// Check if all notes are requested
	allParam := c.Query("all")
	if allParam == "true" {
		limit = 0 // No limit
	} else if limitParam := c.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var notes []*models.Note
	var err error

	if categoryID != "" {
		notes, err = h.repo.GetByCategoryID(categoryID)
	} else {
		notes, err = h.repo.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notes"})
		return
	}

	// Apply limit if needed
	if limit > 0 && len(notes) > limit {
		notes = notes[:limit]
	}

	// Decrypt sensitive data
	var decryptedNotes []*models.Note
	for _, note := range notes {
		decryptedSubject, err := utils.Decrypt(note.Subject)
		if err != nil {
			// Log error but continue with other notes
			fmt.Printf("Error decrypting subject for note %s: %v\n", note.ID, err)
			continue
		}
		note.Subject = decryptedSubject

		decryptedContent, err := utils.Decrypt(note.Content)
		if err != nil {
			fmt.Printf("Error decrypting content for note %s: %v\n", note.ID, err)
			continue
		}
		note.Content = decryptedContent

		decryptedTags, err := utils.Decrypt(note.Tags)
		if err != nil {
			fmt.Printf("Error decrypting tags for note %s: %v\n", note.ID, err)
			continue
		}
		note.Tags = decryptedTags

		// Add successfully decrypted note to the result
		decryptedNotes = append(decryptedNotes, note)
	}

	// Log the activity
	if h.activityLogger != nil {
		h.activityLogger.LogActivity(c, "read", "note", 0, "Retrieved notes")
	}

	c.JSON(http.StatusOK, decryptedNotes)
}

// UpdateNote updates a note by ID
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if note exists
	_, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	note.ID = id

	// Encrypt sensitive data
	encryptedSubject, err := utils.Encrypt(note.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt subject"})
		return
	}
	note.Subject = encryptedSubject

	encryptedContent, err := utils.Encrypt(note.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt content"})
		return
	}
	note.Content = encryptedContent

	encryptedTags, err := utils.Encrypt(note.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt tags"})
		return
	}
	note.Tags = encryptedTags

	if err := h.repo.Update(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	// Log the activity
	if h.activityLogger != nil {
		noteID, _ := strconv.Atoi(note.ID)
		decryptedSubject, _ := utils.Decrypt(note.Subject)
		h.activityLogger.LogActivity(c, "update", "note", noteID, "Updated note: "+decryptedSubject)
	}

	// Decrypt for response
	decryptedSubject, err := utils.Decrypt(note.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt subject"})
		return
	}
	note.Subject = decryptedSubject

	decryptedContent, err := utils.Decrypt(note.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt content"})
		return
	}
	note.Content = decryptedContent

	decryptedTags, err := utils.Decrypt(note.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt tags"})
		return
	}
	note.Tags = decryptedTags

	c.JSON(http.StatusOK, note)
}

// DeleteNote deletes a note by ID
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	// Check if note exists
	_, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	// Log the activity
	if h.activityLogger != nil {
		noteID, _ := strconv.Atoi(id)
		h.activityLogger.LogActivity(c, "delete", "note", noteID, "Deleted note with ID: "+id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
