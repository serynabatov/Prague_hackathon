package password

import (
	"crypto/rand"
	"encoding/base64"
)

// Generates a base64 random password of the specified byte length
func GenerateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
