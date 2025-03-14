package handlers

import (
	"net/http"
	"personal-notes-with-go/models"
	"personal-notes-with-go/repositories"
	"personal-notes-with-go/utils"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	repo repositories.NoteRepository
}

func NewNoteHandler(repo repositories.NoteRepository) *NoteHandler {
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

	if err := h.repo.Create(&note); err != nil {
		utils.HandleInternalServerError(c, err, "create note")
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	priority := c.Query("priority")
	categoryID := c.Query("category_id")

	notes, err := h.repo.GetAll(priority, categoryID)
	if err != nil {
		utils.HandleInternalServerError(c, err, "get notes")
		return
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
