package handlers

import "sync"

type InfoMap struct {
	sync.RWMutex
	keyCount int
	internal map[int]string
}

func NewInfoMap() *InfoMap {
	return &InfoMap{
		internal: make(map[int]string),
	}
}

func (m *InfoMap) Load(key int) string {
	m.RLock()
	defer m.RUnlock()

	result, ok := m.internal[key]

	var hashedPassword string
	if ok {
		hashedPassword = result
	}
	return hashedPassword
}

func (m *InfoMap) GetKey() int {
	m.Lock()
	defer m.Unlock()

	m.keyCount++

	return m.keyCount
}

func (m *InfoMap) Store(key int, hashedPassword string) {
	m.Lock()
	defer m.Unlock()

	m.internal[key] = hashedPassword
}
