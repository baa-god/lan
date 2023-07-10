package typ

import "sync"

type SyncMap[K, V any] struct {
	sync.Map
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	return m.Map.Load(key)
}
