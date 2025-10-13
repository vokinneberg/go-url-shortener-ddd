// Package repository contains persistence abstractions and in-memory implementations.
package repository

import "errors"

// KeyFunc extracts the key of type K from an entity of type T.
type KeyFunc[T any, K comparable] func(*T) K

// InMemoryRepository stores entities of type T in process memory, keyed by K.
// It is intended for testing, demos, and local development.
type InMemoryRepository[T any, K comparable] struct {
	store map[K]*T
	keyFn KeyFunc[T, K]
}

// NewInMemoryRepository creates an empty in-memory repository that derives keys using keyFn.
func NewInMemoryRepository[T any, K comparable](keyFn KeyFunc[T, K]) *InMemoryRepository[T, K] {
	return &InMemoryRepository[T, K]{
		store: make(map[K]*T),
		keyFn: keyFn,
	}
}

// Save upserts the entity in the repository based on its derived key.
func (r *InMemoryRepository[T, K]) Save(entity *T) error {
	r.store[r.keyFn(entity)] = entity
	return nil
}

// Find returns the entity by id or an error if it doesn't exist.
func (r *InMemoryRepository[T, K]) Get(id K) (*T, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return entity, nil
}

// Delete removes the entity from the repository.
func (r *InMemoryRepository[T, K]) Delete(id K) error {
	delete(r.store, id)
	return nil
}
