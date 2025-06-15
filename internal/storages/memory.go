package storages

import "sync"

type MemoryStorage[K comparable, V any] struct {
	Mu    sync.RWMutex
	Store map[K]V
}

func NewMemoryStorage[K comparable, V any]() *MemoryStorage[K, V] {
	return &MemoryStorage[K, V]{
		Store: make(map[K]V),
	}
}
