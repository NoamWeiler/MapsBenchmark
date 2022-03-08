package internal_mutex_map

import (
	"sync"
)

type insideStruct struct {
	s     string
	mutex sync.RWMutex
}

type InternalRWMutexMap struct {
	dbmap map[string]*insideStruct
	sync.RWMutex
}

func NewMap() *InternalRWMutexMap {
	return &InternalRWMutexMap{dbmap: make(map[string]*insideStruct, 100)}
}

func (m *InternalRWMutexMap) Add(key, value string) {
	m.Lock()
	defer m.Unlock()
	m.dbmap[key] = &insideStruct{s: value}
}

func (m *InternalRWMutexMap) Get(key string) string {
	m.Lock()
	defer m.Unlock()
	return m.dbmap[key].s

}

func (m *InternalRWMutexMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.dbmap[key]; !ok {
		return
	}
	delete(m.dbmap, key)
}
