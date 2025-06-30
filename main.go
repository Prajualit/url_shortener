package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string `json:"id"`
	OriginalURL  string `json:"original_url"`
	ShortenedURL string `json:"shortened_url"`
	CreatedAt    string `json:"created_at"`
}

var urlDB = make(map[string]URL)

func generateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	shortURL := fmt.Sprintf("%x", hasher.Sum(nil))[:5]
	fmt.Println("Generated short URL:", shortURL)
	return shortURL
}

func createUrl(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortenedURL: shortURL,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}
	return urlDB[id].ShortenedURL
}

func getURL(id string) (URL, error) {
	if url, ok := urlDB[id]; ok {
		return url, nil
	}
	return URL{}, fmt.Errorf("URL with ID %s not found", id)
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL := "http://localhost:5000/redirect/" + createUrl(data.URL)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"shortened_url": shortURL,
	}
	w.WriteHeader(http.StatusCreated)
	jsonErr := json.NewEncoder(w).Encode(response)
	if jsonErr != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Listening on port 5000\n")
}

func main() {
	fmt.Println("started")

	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	} else {
		fmt.Println("Server running on port 5000")
	}
}
