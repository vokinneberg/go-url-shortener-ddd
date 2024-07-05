package url

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/vokinneberg/go-url-shortener-ddd/domain"
)

type Reader interface {
	Find(id string) (*domain.URL, error)
}

type Writer interface {
	Save(url *domain.URL) error
}

type ReaderWriter interface {
	Reader
	Writer
}

type URLService struct {
	repo ReaderWriter
}

func NewURLService(repo ReaderWriter) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (h *URLService) Shorten(original string) (*domain.URL, error) {
	hash := sha1.New()
	hash.Write([]byte(original))
	short := hex.EncodeToString(hash.Sum(nil))[:8]
	url := domain.NewURL(short, original)
	if err := h.repo.Save(url); err != nil {
		return nil, err
	}
	return url, nil
}

func (h *URLService) Find(id string) (*domain.URL, error) {
	url, err := h.repo.Find(id)
	if err != nil {
		return nil, err
	}
	return url, nil
}
