package url

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/vokinneberg/go-url-shortener-ddd/domain"
)

type URLService struct {
}

func NewURLService() *URLService {
	return &URLService{}
}

func (h *URLService) Shorten(original string) (*domain.URL, error) {
	hash := sha1.New()
	hash.Write([]byte(original))
	short := hex.EncodeToString(hash.Sum(nil))[:8]
	return domain.NewURL(short, original), nil
}

func (h *URLService) Find(id string) (*domain.URL, error) {
	return nil, nil
}
