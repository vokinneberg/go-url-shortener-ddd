package main

import (
	"log"
	"net/http"

	"github.com/vokinneberg/go-url-shortener-ddd/internal/api"
	"github.com/vokinneberg/go-url-shortener-ddd/internal/repository"
	"github.com/vokinneberg/go-url-shortener-ddd/url"
)

func main() {
	repo := repository.NewInMemoryURLRepository()
	urlService := url.NewURLService(repo)
	httpHandler := api.NewHandler(urlService)
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}
