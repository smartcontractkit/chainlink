package pipeline

import (
	"crypto/sha256"
	stdErrors "errors"
	"fmt"
	"sync"

	"github.com/goccy/go-json"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
)

type BridgeConnManager interface {
	GetObservation(bridge bridges.BridgeType, requestData MapParam) ([]byte, error)
}

var (
	ErrBridgeObservationNotFound = stdErrors.New("bridge observation not found")
)

type bridgeConnManager struct {
	mu    sync.RWMutex
	cache map[[32]byte][]byte
}

var defaultBridgeConnManager BridgeConnManager = &bridgeConnManager{
	cache: make(map[[32]byte][]byte),
}

func NewBridgeConnManager() BridgeConnManager {
	return defaultBridgeConnManager
}

// PutObservation writes an observation to the in-memory cache.
func (m *bridgeConnManager) PutObservation(bridge bridges.BridgeType, requestData MapParam, observation []byte) error {
	key, err := bridgeObservationCacheKey(bridge, requestData)
	if err != nil {
		return err
	}
	payload := make([]byte, len(observation))
	copy(payload, observation)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache[key] = payload
	return nil
}

func (m *bridgeConnManager) GetObservation(bridge bridges.BridgeType, requestData MapParam) ([]byte, error) {
	key, err := bridgeObservationCacheKey(bridge, requestData)
	if err != nil {
		return nil, err
	}
	m.mu.RLock()
	entry, ok := m.cache[key]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w for bridge %q", ErrBridgeObservationNotFound, bridge.Name.String())
	}
	payload := make([]byte, len(entry))
	copy(payload, entry)
	return payload, nil
}

func bridgeObservationCacheKey(bridge bridges.BridgeType, requestData MapParam) ([32]byte, error) {
	lookupBytes, err := json.Marshal(requestData)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal bridge lookup payload: %w", err)
	}
	b := make([]byte, 0, len(bridge.Name.String())+len(lookupBytes))
	b = append(b, bridge.Name.String()...)
	b = append(b, lookupBytes...)
	return sha256.Sum256(b), nil
}
