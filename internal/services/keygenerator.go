package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type KeyGenerator struct{}

var kg *KeyGenerator

func NewKeyGenerator() {
	kg = &KeyGenerator{}
}

func GetKeyGenerator() *KeyGenerator {
	return kg
}

func (kg *KeyGenerator) CreateUserKey() (string, error) {
	key, err := generateSecureKey(6)

	if err != nil {
		fmt.Printf("Error in generating key %s", err)
		return "", err
	}

	return key, nil
}

func (kg *KeyGenerator) CreateGameKey() (string, error) {
	key, err := generateSecureKey(5)

	if err != nil {
		fmt.Printf("Error in generating key %s", err)
		return "", err
	}

	return key, nil
}

func generateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Printf("Error in generating secure key %s", err)
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
