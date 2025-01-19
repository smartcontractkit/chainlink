package ethkey

import (
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type State struct {
	ID            int32
	Address       types.EIP55Address
	EVMChainID    big.Big
	Disabled      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	lastUsed      time.Time
	ResourceMutex ResourceMutex // ResourceMutex is an internal field and ought not be persisted to the database.
	// Its main usage is to verify that the same key is not used for both TXMv1 and TXMv2 (usage in both TXMs will cause nonce drift and will lead to missing transactions). This functionality should be removed after we completely switch to TXMv2
}

type ResourceMutex struct {
	mu          sync.Mutex
	activeCount map[ServiceType]int // Tracks active users per service type
}
type ServiceType int

const (
	TXMv1 ServiceType = iota
	TXMv2
)

func (s State) KeyID() string {
	return s.Address.Hex()
}

// lastUsed is an internal field and ought not be persisted to the database or
// exposed outside of the application
func (s State) LastUsed() time.Time {
	return s.lastUsed
}

func (s *State) WasUsed() {
	s.lastUsed = time.Now()
}

// TryLock attempts to lock the resource for the specified service type.
// It returns an error if the resource is locked by a different service type.
func (rm *ResourceMutex) TryLock(serviceType ServiceType) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Check if other service types are using the resource
	for otherServiceType, count := range rm.activeCount {
		if otherServiceType != serviceType && count > 0 {
			return errors.New("resource is locked by another service type")
		}
	}

	// Increment active count for the current service type
	rm.activeCount[serviceType]++
	return nil
}

// Unlock releases the lock for the service type
func (rm *ResourceMutex) Unlock(serviceType ServiceType) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Check if the service type has an active lock
	if rm.activeCount[serviceType] == 0 {
		return errors.New("no active lock for this service type")
	}

	// Decrement active count for the service type
	rm.activeCount[serviceType]--
	if rm.activeCount[serviceType] == 0 {
		delete(rm.activeCount, serviceType)
	}
	return nil
}

// IsLocked checks if the resource is locked by any service or a specific service type.
func (rm *ResourceMutex) IsLocked(serviceType ServiceType) (bool, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Check if the resource is locked by the given service type
	if count, exists := rm.activeCount[serviceType]; exists && count > 0 {
		return true, nil
	}
	return false, nil
}
