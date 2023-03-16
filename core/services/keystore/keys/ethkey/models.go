package ethkey

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type KeyState[ADDR any, ID any, META any] struct {
	ID int32
	//Address    EIP55Address
	Address ADDR
	// EVMChainID utils.Big
	ChainID ID `db:"evm_chain_id"`
	// NextNonce is used for convenience and rendering in UI but the source of
	// truth is always the DB
	// NextNonce int64
	NextMetadata META `db:"next_nonce"`
	Disabled     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	lastUsed     time.Time
	// EVMChainID   utils.Big
	// NextNonce    int64
}

type State KeyState[EIP55Address, utils.Big, int64]

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
