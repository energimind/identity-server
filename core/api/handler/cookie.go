package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"

	"github.com/gin-gonic/gin"
)

const (
	cookieName   = "jwt"
	cookieSecret = "secret-for-my-jwt-cookie"
)

type cookieSecurityContext struct {
	domain string
	secure bool
}

func getCookieSecurityContext(c *gin.Context) cookieSecurityContext {
	domain, _, _ := net.SplitHostPort(c.Request.Host)

	return cookieSecurityContext{
		domain: domain,
		secure: c.Request.URL.Scheme == "https",
	}
}

func encryptCookie(value string) string {
	block, err := aes.NewCipher([]byte(cookieSecret))
	if err != nil {
		panic(err)
	}

	plaintext := []byte(value)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decryptCookie(cryptoText string) (string, error) {
	block, err := aes.NewCipher([]byte(cookieSecret))
	if err != nil {
		return "", fmt.Errorf("new cipher error: %w", err)
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", fmt.Errorf("cookie decode error: %w", err)
	}

	if len(ciphertext) < aes.BlockSize {
		//nolint:goerr113
		return "", fmt.Errorf("ciphertext too short: %d", len(ciphertext))
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
