package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func main() {
	shortener := NewURLShortener()
	r := mux.NewRouter()

	r.HandleFunc("/shorten", shortenHandler(shortener)).Methods("POST")
	r.PathPrefix("/{short_url}").HandlerFunc(redirectHandler(shortener))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func shortenHandler(s *URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		shortID, err := s.Shorten(req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := ShortenResponse{
			ShortURL:    shortID,
			OriginalURL: req.URL,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func redirectHandler(s *URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortID := vars["short_url"]
		orig, err := s.GetOriginal(shortID)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, orig, http.StatusFound)
	}
}
