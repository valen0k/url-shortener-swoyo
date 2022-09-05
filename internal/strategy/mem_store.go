package strategy

import (
	"log"
	"sync"
)

type MemStore struct {
	mem struct {
		sync.RWMutex
		memory map[string]string
	}
}

func NewMemStore() *MemStore {
	mem := MemStore{}
	mem.mem.memory = make(map[string]string)

	return &mem
}

func (m *MemStore) Set(key, val string) error {
	m.mem.Lock()
	defer m.mem.Unlock()
	m.mem.memory[key] = val
	log.Println("recorded in memory")
	return nil
}

func (m *MemStore) Get(key string) (string, bool) {
	m.mem.RLock()
	defer m.mem.RUnlock()

	value, ok := m.mem.memory[key]
	return value, ok
}
