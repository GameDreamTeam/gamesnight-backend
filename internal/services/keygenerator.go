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
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	key, err := generateSecureKey(8, letters)

	if err != nil {
		fmt.Printf("Error in generating key %s", err)
		return "", err
	}

	return key, nil
}

func (kg *KeyGenerator) CreateGameKey() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	key, err := generateSecureKey(4, letters)

	if err != nil {
		fmt.Printf("Error in generating key %s", err)
		return "", err
	}

	return key, nil
}

func generateSecureKey(length int, letters string) (string, error) {
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
