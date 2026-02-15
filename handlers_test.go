package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestShortenHandler(t *testing.T) {
	shortener := NewURLShortener()
	r := mux.NewRouter()
	r.HandleFunc("/shorten", shortenHandler(shortener)).Methods("POST")

	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"valid URL", `{"url": "https://example.com"}`, http.StatusOK},
		{"invalid JSON", "{invalid}", http.StatusBadRequest},
		{"invalid URL", `{"url": "invalid"}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/shorten", bytes.NewReader([]byte(tt.url)))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	shortener := NewURLShortener()
	_, _ = shortener.Shorten("https://example.com")
	r := mux.NewRouter()
	r.PathPrefix("/{short_url}").HandlerFunc(redirectHandler(shortener))

	tests := []struct {
		name       string
		shortID    string
		statusCode int
		location   string
	}{
		{"valid short", "validshortid", http.StatusOK, "/validshortid"}, // assume exists from shorten
		{"not found", "nonexistent", http.StatusNotFound, ""},
	}

	// First create a short ID
	testShortID, _ := shortener.Shorten("https://target.com")
	tests[0].shortID = testShortID
	tests[0].statusCode = http.StatusFound
	tests[0].location = "https://target.com"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+tt.shortID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected %d, got %d", tt.statusCode, w.Code)
			}
			if tt.statusCode == http.StatusFound && w.Header().Get("Location") != tt.location {
				t.Errorf("expected Location %s, got %s", tt.location, w.Header().Get("Location"))
			}
		})
	}
}
