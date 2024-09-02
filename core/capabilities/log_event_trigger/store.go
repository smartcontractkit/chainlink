package log_event_trigger

import (
	"fmt"
	"sync"
)

type NewCapabilityFn[T any, Resp any] func() (T, chan Resp)

// Interface of the capabilities store
type CapabilitiesStore[T any, Resp any] interface {
	Read(capabilityID string) (value T, ok bool)
	ReadAll() (values map[string]T)
	Write(capabilityID string, value T)
	InsertIfNotExists(capabilityID string, fn NewCapabilityFn[T, Resp]) (chan Resp, error)
	Delete(capabilityID string)
}

// Implementation for the CapabilitiesStore interface
type capabilitiesStore[T any, Resp any] struct {
	mu           sync.RWMutex
	capabilities map[string]T
}

var _ CapabilitiesStore[string, string] = (CapabilitiesStore[string, string])(nil)

// Constructor for capabilitiesStore struct implementing CapabilitiesStore interface
func NewCapabilitiesStore[T any, Resp any]() CapabilitiesStore[T, Resp] {
	return &capabilitiesStore[T, Resp]{
		capabilities: map[string]T{},
	}
}

func (cs *capabilitiesStore[T, Resp]) Read(capabilityID string) (value T, ok bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	trigger, ok := cs.capabilities[capabilityID]
	return trigger, ok
}

func (cs *capabilitiesStore[T, Resp]) ReadAll() (values map[string]T) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.capabilities
}

func (cs *capabilitiesStore[T, Resp]) Write(capabilityID string, value T) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.capabilities[capabilityID] = value
}

func (cs *capabilitiesStore[T, Resp]) InsertIfNotExists(capabilityID string, fn NewCapabilityFn[T, Resp]) (chan Resp, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if _, ok := cs.capabilities[capabilityID]; ok {
		return nil, fmt.Errorf("capabilityID %v already exists", capabilityID)
	}
	value, respCh := fn()
	cs.capabilities[capabilityID] = value
	return respCh, nil
}

func (cs *capabilitiesStore[T, Resp]) Delete(capabilityID string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.capabilities, capabilityID)
}
