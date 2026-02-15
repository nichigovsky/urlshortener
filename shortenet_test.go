package main

import (
	"testing"
)

func TestURLShortener_Shorten(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"валидный HTTP URL", "http://example.com", false},
		{"валидный HTTPS URL", "https://google.com/search?q=test", false},
		{"невалидный URL", "not-a-url", true},
		{"пустая строка", "", true},
		{"без схемы", "example.com", true},
		{"ftp URL", "ftp://example.com", true},
	}

	shortener := NewURLShortener()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := shortener.Shorten(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(shortID) != 8 {
				t.Errorf("Shorten() shortID length = %d, want 8", len(shortID))
			}
		})
	}
}

func TestURLShortener_GetOriginal(t *testing.T) {
	shortener := NewURLShortener()
	orig := "https://example.com"
	shortID, _ := shortener.Shorten(orig)

	tests := []struct {
		name    string
		shortID string
		want    string
		wantErr bool
	}{
		{"существующий", shortID, orig, false},
		{"несуществующий", "nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shortener.GetOriginal(tt.shortID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOriginal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetOriginal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateShortID(t *testing.T) {
	id1 := generateShortID()
	id2 := generateShortID()
	if id1 == id2 {
		t.Error("IDs should be different")
	}
	if len(id1) != 8 || len(id2) != 8 {
		t.Errorf("ID length should be 8")
	}
}

func Test_isValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"valid http", "http://example.com", true},
		{"valid https", "https://example.com/path?query=1", true},
		{"invalid scheme", "ftp://example.com", false},
		{"no scheme", "example.com", false},
		{"invalid", "not-url", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURL(tt.url); got != tt.want {
				t.Errorf("isValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
