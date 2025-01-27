package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	n := 32
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}

	encodedStr := hex.EncodeToString(b)
	return encodedStr, nil
}
