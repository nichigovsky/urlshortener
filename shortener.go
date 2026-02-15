package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sync"
)

type URLShortener struct {
    urls map[string]string
    mu   sync.RWMutex
}



func NewURLShortener() *URLShortener {
    return &URLShortener{
        urls: make(map[string]string),
    }
}

// Shorten создает короткий идентификатор для URL
func (us *URLShortener) Shorten(originalURL string) (string, error) {
	if !isValidURL(originalURL) {
		return "", fmt.Errorf("Invalid url: %s", originalURL)
	}

	shortID := generateShortID()

	us.mu.Lock()
	defer us.mu.Unlock()

	us.urls[shortID] = originalURL
	return shortID, nil
}

// GetOriginal возвращает оригинальный URL по короткому ID
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	original, exists := us.urls[shortID]
	if !exists {
		return "", errors.New("short URL not found")
	}
	return original, nil
}

// generateShortID генерирует случайный короткий идентификатор
func generateShortID() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	s := base64.URLEncoding.EncodeToString(bytes)
	return s[:8]
}

// isValidURL проверяет корректность URL
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}