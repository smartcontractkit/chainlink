package ethkey

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type State struct {
	ID         int32
	Address    EIP55Address
	EVMChainID utils.Big
	// NextNonce is used for convenience and rendering in UI but the source of
	// truth is always the DB
	NextNonce int64
	Disabled  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	lastUsed  time.Time
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
