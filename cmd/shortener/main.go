package main

import (
	"log"
	"net/http"

	"github.com/vokinneberg/go-url-shortener-ddd/internal/api"
)

func main() {
	httpHandler := api.NewHandler()
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}
