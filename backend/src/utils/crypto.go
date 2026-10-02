package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func Hash(str string) string {
	bytes := sha256.Sum256([]byte(str))
	h := base64.RawURLEncoding.EncodeToString(bytes[:])
	return h
}

func GenerateId() (string, error) {
	bytes := [24]byte{}
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}

	str := base64.RawURLEncoding.EncodeToString(bytes[:])

	return str, nil
}

func GenerateBase64(n int) (string, error) {
	bytes := make([]byte, (n/4)*3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	str := base64.RawURLEncoding.EncodeToString(bytes)

	return str, nil
}
