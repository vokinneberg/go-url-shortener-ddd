package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vokinneberg/go-url-shortener-ddd/internal/api"
	"github.com/vokinneberg/go-url-shortener-ddd/internal/config"
	"github.com/vokinneberg/go-url-shortener-ddd/internal/repository"
	"github.com/vokinneberg/go-url-shortener-ddd/internal/url"
)

var _ url.URLReaderWriter = (*repository.InMemoryRepository[url.URL, string])(nil)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	repo := repository.NewInMemoryRepository(func(url *url.URL) string {
		return url.ID
	})
	urlService := url.NewURLService(repo, func(original string) (string, error) {
		hash := sha1.New()
		hash.Write([]byte(original))
		short := hex.EncodeToString(hash.Sum(nil))[:8]
		return short, nil
	})
	httpHandler := api.NewHandler(urlService)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: httpHandler,
	}

	// Start server in a separate goroutine
	go func() {
		log.Println("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Listen for interrupt/termination signals.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done() // Block until a signal is received.
	log.Println("Shutdown signal received")

	// Give outstanding requests some time to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
		// Force close if needed
		if cerr := srv.Close(); cerr != nil {
			log.Printf("Forced close error: %v", cerr)
		}
	} else {
		log.Println("Server shutdown complete")
	}
}
