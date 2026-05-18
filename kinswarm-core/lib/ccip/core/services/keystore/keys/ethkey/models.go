package ethkey

import (
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
