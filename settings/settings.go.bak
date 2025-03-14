package settings

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type Settings struct {
	EncryptionKey string `json:"encryption_key"`
}

const (
	settingsFile = "settings.json"
	keyLength    = 32 // Length of the encryption key in bytes
)

// LoadSettings loads settings from the settings.json file
// If the file doesn't exist or the encryption key is not set,
// it will generate a new key and save it
func LoadSettings() (*Settings, error) {
	var settings Settings

	// Try to read existing settings
	data, err := os.ReadFile(settingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, generate new settings
			return generateAndSaveSettings()
		}
		return nil, err
	}

	// Parse existing settings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		return nil, err
	}

	// If encryption key is not set, generate new settings
	if settings.EncryptionKey == "" {
		return generateAndSaveSettings()
	}

	return &settings, nil
}

// generateAndSaveSettings creates new settings with a generated encryption key
// and saves them to the settings file
func generateAndSaveSettings() (*Settings, error) {
	// Generate new encryption key
	key := generateEncryptionKey()

	// Create settings with the new key
	settings := &Settings{
		EncryptionKey: base64.StdEncoding.EncodeToString(key),
	}

	// Save to file
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(settingsFile, data, 0600)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// generateEncryptionKey creates a new random encryption key
func generateEncryptionKey() []byte {
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, keyLength)
	rand.Read(key)
	return key
}

// GetEncryptionKey returns the base64-decoded encryption key
func (s *Settings) GetEncryptionKey() ([]byte, error) {
	return base64.StdEncoding.DecodeString(s.EncryptionKey)
}
