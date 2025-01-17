package ethkey

import (
	"errors"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type State struct {
	ID         int32
	Address    types.EIP55Address
	EVMChainID big.Big
	Disabled   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	lastUsed   time.Time
	tags       map[string]int // tags is an internal field and ought not be persisted to the database.
	// Its main usage is to verify that the same key is not used for both TXMv1 and TXMv2(aka primary and secondary transmitter)
	// This functionality should be removed after we completely switch to TXMv2
}

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

// Tag Adds a process to the list of processes that have are using
func (s *State) Tag(label string) error {
	if label == "" {
		return errors.New("cannot add usage: label string is empty")
	}
	label = strings.ToLower(label)

	if _, exists := s.tags[label]; exists {
		s.tags[label]++
	} else {
		s.tags[label] = 1
	}
	return nil
}

// Untag removes a process from the list of processes that have used the key
func (s *State) Untag(label string) error {
	if label == "" {
		return errors.New("cannot remove usage: label string is empty")
	}
	label = strings.ToLower(label)

	if count, exists := s.tags[label]; exists {
		if count > 1 {
			s.tags[label]--
		} else {
			delete(s.tags, label)
		}
		return nil
	}
	return nil
}

// HasTag checks if the key is used by the given process
func (s *State) HasTag(label string) (bool, error) {
	if label == "" {
		return false, errors.New("cannot check usage: label string is empty")
	}
	label = strings.ToLower(label)
	_, exists := s.tags[label]
	return exists, nil
}
