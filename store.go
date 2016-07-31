package main

import (
	"sync"
	"time"
)

// Store TODO
type Store interface {
	Store(key []byte, data []byte, expiration time.Time, publisher bool) error
	Retrieve(key []byte) ([]byte, bool)
	Delete(key []byte)
	Init()
	GetAllKeysForRefresh() []string
	ExpireKeys()
}

// MemoryStore TODO
type MemoryStore struct {
	mutex      *sync.Mutex
	data       map[string][]byte
	refreshMap map[string]time.Time
	expireMap  map[string]time.Time
}

// GetAllKeysForRefresh TODO
func (ms *MemoryStore) GetAllKeysForRefresh() []string {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	var keys []string
	for k := range ms.data {
		if time.Now().After(ms.refreshMap[k]) {
			keys = append(keys, k)
		}
	}
	return keys
}

// ExpireKeys TODO
func (ms *MemoryStore) ExpireKeys() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	for k, v := range ms.expireMap {
		if time.Now().After(v) {
			delete(ms.refreshMap, k)
			delete(ms.expireMap, k)
			delete(ms.data, k)
		}
	}
}

// Init TODO
func (ms *MemoryStore) Init() {
	ms.data = make(map[string][]byte)
	ms.mutex = &sync.Mutex{}
	ms.refreshMap = make(map[string]time.Time)
	ms.expireMap = make(map[string]time.Time)
}

// Store TODO
func (ms *MemoryStore) Store(key []byte, data []byte, expiration time.Time, publisher bool) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.refreshMap[string(key)] = time.Now().Add(time.Hour * 1)
	ms.expireMap[string(key)] = expiration
	ms.data[string(key)] = data
	return nil
}

// Retrieve TODO
func (ms *MemoryStore) Retrieve(key []byte) ([]byte, bool) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	data, exists := ms.data[string(key)]
	return data, exists
}

// Delete TODO
func (ms *MemoryStore) Delete(key []byte) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	delete(ms.data, string(key))
}