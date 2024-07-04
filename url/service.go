package url

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/vokinneberg/go-url-shortener-ddd/domain"
)

type URLFetcher interface {
	FindByID(id string) (*domain.URL, error)
}

type URLSaver interface {
	Save(url *domain.URL) error
}

type URLRepo interface {
	URLFetcher
	URLSaver
}

type URLService struct {
	repo URLRepo
}

func NewURLService(repo URLRepo) *URLService {
	return &URLService{repo: repo}
}

func (h *URLService) Shorten(original string) (*domain.URL, error) {
	hash := sha1.New()
	hash.Write([]byte(original))
	short := hex.EncodeToString(hash.Sum(nil))[:8]
	return domain.NewURL(short, original), nil
}

func (h *URLService) Find(id string) (*domain.URL, error) {
	url, err := h.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return url, nil
}
