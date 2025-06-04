package fakes

import (
	"sync"
)

type cronStore struct {
	mu       sync.RWMutex
	triggers map[string]cronTrigger
}

type CronStore interface {
	Read(triggerID string) (value cronTrigger, ok bool)
	ReadAll() (values map[string]cronTrigger)
	Write(triggerID string, value cronTrigger)
	Delete(triggerID string)
}

var _ CronStore = (CronStore)(nil)

func NewCronStore() *cronStore {
	return &cronStore{
		triggers: map[string]cronTrigger{},
	}
}

func (cs *cronStore) Read(triggerID string) (value cronTrigger, ok bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	trigger, ok := cs.triggers[triggerID]
	return trigger, ok
}

func (cs *cronStore) ReadAll() (values map[string]cronTrigger) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	tCopy := map[string]cronTrigger{}
	for key, value := range cs.triggers {
		tCopy[key] = value
	}
	return tCopy
}

func (cs *cronStore) Write(triggerID string, value cronTrigger) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.triggers[triggerID] = value
}

func (cs *cronStore) Delete(triggerID string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.triggers, triggerID)
}
