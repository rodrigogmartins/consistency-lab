package store

import (
	"sync"
	"time"
)

type Item struct {
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	Version   int64     `json:"version"` // monotonic per-node logical clock (good enough for demo)
	UpdatedAt time.Time `json:"updatedAt"`
}

type Store struct {
	mu    sync.RWMutex
	items map[string]Item
	clock int64
	node  string
}

func New(node string) *Store {
	return &Store{
		items: make(map[string]Item),
		node:  node,
	}
}

func (s *Store) Put(id, value string) Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clock++
	it := Item{
		ID:        id,
		Value:     value,
		Version:   s.clock,
		UpdatedAt: time.Now(),
	}
	s.items[id] = it
	return it
}

// ApplyReplica applies an incoming replica update.
// Conflict rule: last-writer-wins by Version; since versions are per-node, this is "good enough" for demo.
// (You can upgrade to (node,version) or HLC later.)
func (s *Store) ApplyReplica(it Item) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	cur, ok := s.items[it.ID]
	if !ok || it.Version >= cur.Version {
		s.items[it.ID] = it
		return true
	}
	return false
}

func (s *Store) Get(id string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	return it, ok
}
