package utils

import (
	"fmt"
	"github.com/jkomyno/nanoid"
)

func GenerateUniqueToken () (string, error) {
	token, err := nanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
	if err != nil {
		return "", fmt.Errorf("nanoid.Generate: %v", err)
	}
	return token, nil
}

func GenerateEncryptionKey () (string, error) {
	key, err := nanoid.Nanoid(32)
	if err != nil {
		return "", fmt.Errorf("nanoid.Nanoid: %v", err)
	}
	return key, nil
}

func GenerateFileName () (string, error) {
	key, err := nanoid.Generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz-" ,30)
	if err != nil {
		return "", fmt.Errorf("nanoid.Nanoid: %v", err)
	}
	return key, nil
}
