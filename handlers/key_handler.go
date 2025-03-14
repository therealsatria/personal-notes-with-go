package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"personal-notes-with-go/utils"

	"github.com/gin-gonic/gin"
)

type KeyGenerateRequest struct {
	Text string `json:"text" binding:"required"`
}

type KeyHandler struct{}

func NewKeyHandler() *KeyHandler {
	return &KeyHandler{}
}

func (h *KeyHandler) GenerateKey(c *gin.Context) {
	var req KeyGenerateRequest
	if err := c.BindJSON(&req); err != nil {
		utils.HandleBadRequestError(c, err)
		return
	}

	if req.Text == "" {
		utils.HandleBadRequestError(c, utils.ErrEmptyInput)
		return
	}

	// Generate a consistent key using SHA-256 and Base64
	hasher := sha256.New()
	hasher.Write([]byte(req.Text))
	hash := hasher.Sum(nil)
	key := base64.StdEncoding.EncodeToString(hash)

	c.JSON(http.StatusOK, gin.H{
		"key": key,
	})
}
