package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	shortURL "github.com/vokinneberg/go-url-shortener-ddd/url"
)

type Handler struct {
	urlService *shortURL.URLService
}

func NewHandler() *http.ServeMux {
	handler := &Handler{
		urlService: shortURL.NewURLService(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handler.shortenURL)
	mux.HandleFunc("GET /{id}", handler.getURL)
	return mux
}

func (h *Handler) shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}

	url, err := url.QueryUnescape(string(body))
	if err != nil {
		http.Error(w, "Unable to decode URL-encoded string", http.StatusBadRequest)
		return
	}

	shortURL, err := h.urlService.Shorten(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, shortURL.ID))); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}

func (h *Handler) getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	shortURL := r.PathValue("id")

	origURL, err := h.urlService.Find(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", origURL.Original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
