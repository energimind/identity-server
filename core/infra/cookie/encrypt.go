package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func encryptCookie(value, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("new cipher error: %w", err)
	}

	plainText := []byte(value)
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]

	if _, rfErr := io.ReadFull(rand.Reader, iv); rfErr != nil {
		return "", fmt.Errorf("read full error: %w", rfErr)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func decryptCookie(value, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("new cipher error: %w", err)
	}

	cipherText, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("cookie decode error: %w", err)
	}

	if len(cipherText) < aes.BlockSize {
		//nolint:goerr113
		return "", fmt.Errorf("cipherText too short: %d", len(cipherText))
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
