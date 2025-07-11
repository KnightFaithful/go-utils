package task

import (
	"errors"
	"sync"
)

type ConcurrentIntMap struct {
	Map  map[interface{}]interface{}
	Lock sync.RWMutex
}

func NewConnManger() *ConcurrentIntMap {
	cm := &ConcurrentIntMap{
		Map: make(map[interface{}]interface{}),
	}
	return cm
}

func (cm *ConcurrentIntMap) Add(id interface{}, value interface{}) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	cm.Map[id] = value
}

func (cm *ConcurrentIntMap) Remove(id interface{}) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	delete(cm.Map, id)
}
func (cm *ConcurrentIntMap) Get(id interface{}) (interface{}, error) {
	cm.Lock.RLock()
	defer cm.Lock.RUnlock()
	conn, ok := cm.Map[id]
	if !ok {
		return "", errors.New("connmanager get conn error ")
	}
	return conn, nil
}
func (cm *ConcurrentIntMap) Len() int {
	return len(cm.Map)
}
func (cm *ConcurrentIntMap) Clean() {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	for key, _ := range cm.Map {
		delete(cm.Map, key)
	}
}
