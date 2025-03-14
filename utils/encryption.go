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
	encryptionKey   []byte
	encryptionValid bool // Flag to track if encryption is valid
)

// InitEncryption initializes the encryption system with the key from settings
func InitEncryption() error {
	// Reset encryption status
	encryptionValid = false

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

	// Validate encryption key by performing a test encryption and decryption
	if err := validateEncryptionKey(); err != nil {
		return err
	}

	// If we reach here, encryption is valid
	encryptionValid = true
	return nil
}

// validateEncryptionKey tests if the encryption key is valid by encrypting and decrypting a test string
func validateEncryptionKey() error {
	testString := "encryption_test"

	// Try to encrypt
	encrypted, err := encryptWithKey(testString, encryptionKey)
	if err != nil {
		return errors.New("encryption key validation failed: " + err.Error())
	}

	// Try to decrypt
	decrypted, err := decryptWithKey(encrypted, encryptionKey)
	if err != nil {
		return errors.New("encryption key validation failed: " + err.Error())
	}

	// Check if decrypted matches original
	if decrypted != testString {
		return errors.New("encryption key validation failed: decrypted text does not match original")
	}

	return nil
}

// IsEncryptionValid returns whether the encryption system is properly initialized and validated
func IsEncryptionValid() bool {
	return encryptionValid
}

// Encrypt encrypts the given text using AES-256 and returns a base64 encoded string
func Encrypt(text string) (string, error) {
	if !encryptionValid {
		return "", errors.New("encryption system not properly initialized")
	}

	return encryptWithKey(text, encryptionKey)
}

// encryptWithKey encrypts text with the provided key
func encryptWithKey(text string, key []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("encryption key not provided")
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
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
	if !encryptionValid {
		return "", errors.New("encryption system not properly initialized")
	}

	return decryptWithKey(encryptedText, encryptionKey)
}

// decryptWithKey decrypts text with the provided key
func decryptWithKey(encryptedText string, key []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("encryption key not provided")
	}

	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
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
