package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"personal-notes-with-go/settings"
)

var (
	encryptionKey []byte
)

// InitEncryption initializes the encryption system with the key from settings
func InitEncryption() error {
	// Load settings
	s, err := settings.LoadSettings()
	if err != nil {
		return err
	}

	// Get the encryption key
	key, err := s.GetEncryptionKey()
	if err != nil {
		return err
	}

	encryptionKey = key
	return nil
}

// Encrypt encrypts the given text using AES-256 and returns a base64 encoded string
func Encrypt(text string) (string, error) {
	if len(encryptionKey) == 0 {
		return "", errors.New("encryption key not initialized")
	}

	// Create cipher block
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given base64 encoded ciphertext
func Decrypt(encryptedText string) (string, error) {
	if len(encryptionKey) == 0 {
		return "", errors.New("encryption key not initialized")
	}

	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Create cipher block
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
