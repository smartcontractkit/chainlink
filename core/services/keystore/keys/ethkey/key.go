package ethkey

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg/datatypes"
)

// NOTE: This model refers to the OLD key and is only used for migrations
//
// Key holds the private key metadata for a given address that is used to unlock
// said key when given a password.
//
// By default, a key is assumed to represent an ethereum account.
type Key struct {
	ID        int32
	Address   EIP55Address
	JSON      datatypes.JSON `json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt *time.Time     `json:"-"`
	// This is the nonce that should be used for the next transaction.
	// Conceptually equivalent to geth's `PendingNonceAt` but more reliable
	// because we have a better view of our own transactions
	// NOTE: Be cautious about using this field, it is provided for convenience
	// only, can go out of date, and should not be relied upon. The source of
	// truth is always the database row for the key.
	NextNonce int64 `json:"-"`
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
