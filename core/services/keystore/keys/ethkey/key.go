package ethkey

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// NOTE: This model refers to the OLD key and is only used for migrations
//
// Key holds the private key metadata for a given address that is used to unlock
// said key when given a password.
//
// By default, a key is assumed to represent an ethereum account.
type Key struct {
	ID        int32
	Address   types.EIP55Address
	JSON      sqlutil.JSON `json:"-"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt *time.Time   `json:"-"`
	// IsFunding marks the address as being used for rescuing the  node and the pending transactions
	// Only one key can be IsFunding=true at a time.
	IsFunding bool
}

// Type returns type of key
func (k Key) Type() string {
	if k.IsFunding {
		return "funding"
	}
	return "sending"
}
