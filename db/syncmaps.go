package db

import (
	"sync"
)

type ConcurrentMap[K comparable, V any] struct {
	sync.RWMutex
	Map map[K]V
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		Map: make(map[K]V),
	}
}

func (cm *ConcurrentMap[K, V]) Del(key K) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.Map, key)
}

func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	cm.Lock()
	defer cm.Unlock()
	val, ok := cm.Map[key]
	return val, ok
}

func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
	cm.Lock()
	defer cm.Unlock()
	cm.Map[key] = value
}
