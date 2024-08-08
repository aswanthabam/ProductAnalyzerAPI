package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateAPIKey generates a random API key and encodes it to base64
func GenerateAPIKey() (string, error) {
	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
