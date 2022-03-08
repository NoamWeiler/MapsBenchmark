package rw_map

import (
	"sync"
)

type RWMap struct {
	dbmap map[string]string
	sync.RWMutex
}

func NewMap() *RWMap {
	return &RWMap{dbmap: make(map[string]string, 100)}
}

func (m *RWMap) Add(key, value string) {
	m.Lock()
	defer m.Unlock()
	m.dbmap[key] = value
}

func (m *RWMap) Get(key string) string {
	m.RLock()
	defer m.RUnlock()
	if v, ok := m.dbmap[key]; !ok {
		return ""
	} else {
		return v
	}
}

func (m *RWMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.dbmap[key]; !ok {
		return
	}
	delete(m.dbmap, key)
}
