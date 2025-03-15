package handlers

import (
	"net/http"
	"personal-notes-with-go/utils"

	"github.com/gin-gonic/gin"
)

// EncryptionHandler handles encryption-related endpoints
type EncryptionHandler struct {
	activityLogger *ActivityLogHandler
}

// NewEncryptionHandler creates a new encryption handler
func NewEncryptionHandler() *EncryptionHandler {
	return &EncryptionHandler{}
}

// SetActivityLogger sets the activity logger for this handler
func (h *EncryptionHandler) SetActivityLogger(logger *ActivityLogHandler) {
	h.activityLogger = logger
}

// GetStatus returns the current status of the encryption system
func (h *EncryptionHandler) GetStatus(c *gin.Context) {
	isValid := utils.IsEncryptionValid()

	// Log activity
	if h.activityLogger != nil {
		status := "valid"
		if !isValid {
			status = "invalid"
		}
		description := "Checked encryption status: " + status
		h.activityLogger.LogActivity(c, "check", "encryption", 0, description)
	}

	c.JSON(http.StatusOK, gin.H{
		"encryption_valid": isValid,
		"message":          getEncryptionStatusMessage(isValid),
	})
}

// Helper function to get a user-friendly message based on encryption status
func getEncryptionStatusMessage(isValid bool) string {
	if isValid {
		return "Encryption system is properly initialized and working correctly."
	}
	return "Encryption system is not properly initialized. Data modification is disabled for security reasons."
}
