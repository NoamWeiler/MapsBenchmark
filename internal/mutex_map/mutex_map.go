package mutex_map

import (
	"sync"
)

type MutexMap struct {
	dbmap map[string]string
	sync.Mutex
}

func NewMap() *MutexMap {
	return &MutexMap{dbmap: make(map[string]string, 100)}
}

func (m *MutexMap) Add(key, value string) {
	m.Lock()
	defer m.Unlock()
	m.dbmap[key] = value
}

func (m *MutexMap) Get(key string) string {
	m.Lock()
	defer m.Unlock()
	return m.dbmap[key]

}

func (m *MutexMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.dbmap[key]; !ok {
		return
	}
	delete(m.dbmap, key)
}
