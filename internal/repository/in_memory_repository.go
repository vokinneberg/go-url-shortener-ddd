package repository

import (
	"errors"

	"github.com/vokinneberg/go-url-shortener-ddd/domain"
)

type InMemoryURLRepository struct {
	store map[string]*domain.URL
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{store: make(map[string]*domain.URL)}
}

func (r *InMemoryURLRepository) Save(url *domain.URL) error {
	r.store[url.ID] = url
	return nil
}

func (r *InMemoryURLRepository) Find(id string) (*domain.URL, error) {
	url, ok := r.store[id]
	if !ok {
		return nil, errors.New("URL not found")
	}
	return url, nil
}
