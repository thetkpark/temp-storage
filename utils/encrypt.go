package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
)

func EncryptFile(file *multipart.File, key string) (*bytes.Buffer, error) {

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, *file); err != nil {
		return nil, fmt.Errorf("io.Copy: %v", err)
	}

	cyphr, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %v", err)
	}

	gcm, err := cipher.NewGCM(cyphr)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("io.ReadFull: %v", err)
	}

	encrypted := gcm.Seal(nonce, nonce, buffer.Bytes(), nil)

	encryptedBuffer := bytes.NewBuffer(encrypted)
	return encryptedBuffer, nil
}

func DecryptFile(file *[]byte, key string) *[]byte {

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	gcmDecrypt, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonceSize := gcmDecrypt.NonceSize()
	if len(*file) < nonceSize {
		fmt.Println(err)
	}
	nonce, encryptedMessage := (*file)[:nonceSize], (*file)[nonceSize:]
	decryptData, err := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		fmt.Println(err)
	}
	return &decryptData
}