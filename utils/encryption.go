package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// ErrInvalidKeyLength is returned when the key length is not 32 bytes.
var ErrInvalidKeyLength = errors.New("key length must be 32 bytes for AES-256")

// loadEncryptionKey loads the encryption key from a file.
func loadEncryptionKey(keyFilePath string) ([]byte, error) {
	key, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encryption key file: %w", err)
	}

	if len(key) != 32 {
		return nil, ErrInvalidKeyLength
	}

	return key, nil
}

// encrypt encrypts the given data using AES-256.
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to create nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts the given data using AES-256.
func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string using AES-256.
func EncryptString(text string, keyFilePath string) (string, error) {
	key, err := loadEncryptionKey(keyFilePath)
	if err != nil {
		return "", err
	}

	ciphertext, err := encrypt([]byte(text), key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString decrypts a string using AES-256.
func DecryptString(encryptedText string, keyFilePath string) (string, error) {
	key, err := loadEncryptionKey(keyFilePath)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
