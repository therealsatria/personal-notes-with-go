package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KeyGenerateRequest struct {
	Text string `json:"text" binding:"required"`
}

type KeyHandler struct {
	activityLogger *ActivityLogHandler
}

func NewKeyHandler() *KeyHandler {
	return &KeyHandler{}
}

// SetActivityLogger sets the activity logger for this handler
func (h *KeyHandler) SetActivityLogger(logger *ActivityLogHandler) {
	h.activityLogger = logger
}

func (h *KeyHandler) GenerateKey(c *gin.Context) {
	var req KeyGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text cannot be empty"})
		return
	}

	// Generate a consistent key using SHA-256 and Base64
	hasher := sha256.New()
	hasher.Write([]byte(req.Text))
	hash := hasher.Sum(nil)
	key := base64.StdEncoding.EncodeToString(hash)

	// Log activity
	if h.activityLogger != nil {
		description := "Generated encryption key"
		h.activityLogger.LogActivity(c, "generate", "key", 0, description)
	}

	c.JSON(http.StatusOK, gin.H{
		"key": key,
	})
}
