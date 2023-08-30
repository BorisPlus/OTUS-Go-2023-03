package transmitter

import (
	"sync"
)

type HashFunc[T Item] func(T) string

type Set[T Item] struct {
	internal map[string]T
	mutex    *sync.RWMutex
	Hash     HashFunc[T]
}

func NewSet[T Item](hashFunc HashFunc[T]) *Set[T] {
	return &Set[T]{
		make(map[string]T),
		&sync.RWMutex{},
		hashFunc,
	}
}

func (s *Set[T]) add(element T) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.internal[s.Hash(element)] = element
}

func (s *Set[T]) remove(element T) { //nolint:unused
	defer s.mutex.RUnlock()
	s.mutex.RLock()
	delete(s.internal, s.Hash(element))
}

func (s *Set[T]) clear() {
	defer s.mutex.RUnlock()
	s.mutex.RLock()
	s.internal = map[string]T{}
}

func (s *Set[T]) has(element T) bool {
	defer s.mutex.RUnlock()
	s.mutex.RLock()
	_, ok := s.internal[s.Hash(element)]
	return ok
}
