package storages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_BasicOperations(t *testing.T) {
	type Key string
	type Value int

	storage := NewMemoryStorage[Key, Value]()

	// Данные для теста
	k1 := Key("key1")
	v1 := Value(42)

	k2 := Key("key2")
	v2 := Value(100)

	// Тест: запись
	storage.Mu.Lock()
	storage.Store[k1] = v1
	storage.Store[k2] = v2
	storage.Mu.Unlock()

	// Тест: чтение
	storage.Mu.RLock()
	val1, ok1 := storage.Store[k1]
	val2, ok2 := storage.Store[k2]
	storage.Mu.RUnlock()

	assert.True(t, ok1, "key1 should exist")
	assert.Equal(t, v1, val1, "key1 value mismatch")

	assert.True(t, ok2, "key2 should exist")
	assert.Equal(t, v2, val2, "key2 value mismatch")

	// Тест: удаление
	storage.Mu.Lock()
	delete(storage.Store, k1)
	storage.Mu.Unlock()

	storage.Mu.RLock()
	_, ok1 = storage.Store[k1]
	storage.Mu.RUnlock()

	assert.False(t, ok1, "key1 should be deleted")
}
