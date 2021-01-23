package main

import (
	"fmt"
	"github.com/jkomyno/nanoid"
)

func generateUniqueToken () (string, error) {
	token, err := nanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
	if err != nil {
		return "", fmt.Errorf("nanoid.Generate: %v", err)
	}
	return token, nil
}
