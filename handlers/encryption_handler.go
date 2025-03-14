package handlers

import (
	"net/http"
	"personal-notes-with-go/utils"

	"github.com/gin-gonic/gin"
)

// EncryptionHandler handles encryption-related endpoints
type EncryptionHandler struct{}

// NewEncryptionHandler creates a new encryption handler
func NewEncryptionHandler() *EncryptionHandler {
	return &EncryptionHandler{}
}

// GetStatus returns the current status of the encryption system
func (h *EncryptionHandler) GetStatus(c *gin.Context) {
	isValid := utils.IsEncryptionValid()

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
