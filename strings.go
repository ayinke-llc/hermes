package hermes

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func IsStringEmpty(s string) bool { return len(strings.TrimSpace(s)) == 0 }

func Random(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
