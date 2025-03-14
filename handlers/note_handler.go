package handlers

import (
	"net/http"
	"personal-notes-with-go/models"
	"personal-notes-with-go/repositories"
	"personal-notes-with-go/settings"
	"personal-notes-with-go/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	repo repositories.NoteRepositoryInterface
}

func NewNoteHandler(repo repositories.NoteRepositoryInterface) *NoteHandler {
	return &NoteHandler{repo: repo}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var note models.Note
	if err := c.BindJSON(&note); err != nil {
		utils.HandleBadRequestError(c, err)
		return
	}

	if note.Subject == "" {
		utils.HandleBadRequestError(c, utils.ErrNoteSubjectEmpty)
		return
	}

	// Clear any ID that might have been sent in the request
	// ID will be generated by the repository
	note.ID = ""

	if err := h.repo.Create(&note); err != nil {
		utils.HandleInternalServerError(c, err, "create note")
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	// Check if search query is provided
	query := c.Query("q")
	// Check if all notes should be returned
	showAllParam := c.Query("all")
	showAll := showAllParam == "true"

	// Load settings to get the notes limit
	s, err := settings.LoadSettings()
	if err != nil {
		utils.HandleInternalServerError(c, err, "load settings")
		return
	}

	// Get limit from settings, default is handled by GetNotesLimit
	limit := s.GetNotesLimit()
	if showAll {
		limit = 0 // No limit
	}

	notes, err := h.repo.GetAll()
	if err != nil {
		utils.HandleInternalServerError(c, err, "get notes")
		return
	}

	// If search query is provided, filter notes
	if query != "" {
		query = strings.ToLower(query)
		var filteredNotes []*models.Note

		for _, note := range notes {
			// Search in subject and content (both already decrypted by repository)
			if strings.Contains(strings.ToLower(note.Subject), query) ||
				strings.Contains(strings.ToLower(note.Content), query) {
				filteredNotes = append(filteredNotes, note)
			}
		}

		notes = filteredNotes
	}

	// Apply limit if needed
	if !showAll && limit > 0 && len(notes) > limit {
		notes = notes[:limit]
	}

	c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := c.BindJSON(&note); err != nil {
		utils.HandleBadRequestError(c, err)
		return
	}

	if note.Subject == "" {
		utils.HandleBadRequestError(c, utils.ErrNoteSubjectEmpty)
		return
	}

	note.ID = id
	if err := h.repo.Update(&note); err != nil {
		utils.HandleInternalServerError(c, err, "update note")
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		utils.HandleInternalServerError(c, err, "delete note")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}
