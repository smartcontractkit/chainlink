package logevent

import (
	"fmt"
	"sync"
)

type RegisterCapabilityFn[T any, Resp any] func() (*T, chan Resp, error)

// Interface of the capabilities store
type CapabilitiesStore[T any, Resp any] interface {
	Read(capabilityID string) (value *T, ok bool)
	ReadAll() (values []*T)
	Write(capabilityID string, value *T)
	InsertIfNotExists(capabilityID string, fn RegisterCapabilityFn[T, Resp]) (chan Resp, error)
	Delete(capabilityID string)
}

// Implementation for the CapabilitiesStore interface
type capabilitiesStore[T any, Resp any] struct {
	mu           sync.RWMutex
	capabilities map[string]*T
}

var _ CapabilitiesStore[string, string] = (CapabilitiesStore[string, string])(nil)

// Constructor for capabilitiesStore struct implementing CapabilitiesStore interface
func NewCapabilitiesStore[T any, Resp any]() CapabilitiesStore[T, Resp] {
	return &capabilitiesStore[T, Resp]{
		capabilities: map[string]*T{},
	}
}

func (cs *capabilitiesStore[T, Resp]) Read(capabilityID string) (value *T, ok bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	trigger, ok := cs.capabilities[capabilityID]
	return trigger, ok
}

func (cs *capabilitiesStore[T, Resp]) ReadAll() (values []*T) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	vals := make([]*T, 0)
	for _, v := range cs.capabilities {
		vals = append(vals, v)
	}
	return vals
}

func (cs *capabilitiesStore[T, Resp]) Write(capabilityID string, value *T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.capabilities[capabilityID] = value
}

func (cs *capabilitiesStore[T, Resp]) InsertIfNotExists(capabilityID string, fn RegisterCapabilityFn[T, Resp]) (chan Resp, error) {
	cs.mu.RLock()
	_, ok := cs.capabilities[capabilityID]
	cs.mu.RUnlock()
	if ok {
		return nil, fmt.Errorf("capabilityID %v already exists", capabilityID)
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	_, ok = cs.capabilities[capabilityID]
	if ok {
		return nil, fmt.Errorf("capabilityID %v already exists", capabilityID)
	}
	value, respCh, err := fn()
	if err != nil {
		return nil, fmt.Errorf("error registering capability: %v", err)
	}
	cs.capabilities[capabilityID] = value
	return respCh, nil
}

func (cs *capabilitiesStore[T, Resp]) Delete(capabilityID string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.capabilities, capabilityID)
}
