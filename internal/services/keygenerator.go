package services

import (
	"crypto/rand"
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
	key, err := generateSecureKey(8)

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
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Printf("Error in generating secure key %s", err)
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes), nil
}
