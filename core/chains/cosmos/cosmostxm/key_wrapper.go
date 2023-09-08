package cosmostxm

import (
	"bytes"
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// KeyWrapper uses a KeystoreAdapter to implement the cosmos-sdk PrivKey interface for a specific key.
type KeyWrapper struct {
	adapter *KeystoreAdapter
	account string
}

var _ cryptotypes.PrivKey = &KeyWrapper{}

func NewKeyWrapper(adapter *KeystoreAdapter, account string) *KeyWrapper {
	return &KeyWrapper{
		adapter: adapter,
		account: account,
	}
}

func (a *KeyWrapper) Bytes() []byte {
	// don't expose the private key.
	return nil
}

func (a *KeyWrapper) Sign(msg []byte) ([]byte, error) {
	return a.adapter.Sign(context.Background(), a.account, msg)
}

func (a *KeyWrapper) PubKey() cryptotypes.PubKey {
	pubKey, err := a.adapter.PubKey(a.account)
	if err != nil {
		// return an empty pubkey if it's not found.
		return &secp256k1.PubKey{Key: []byte{}}
	}
	return pubKey
}

func (a *KeyWrapper) Equals(other cryptotypes.LedgerPrivKey) bool {
	return bytes.Equal(a.PubKey().Bytes(), other.PubKey().Bytes())
}

func (a *KeyWrapper) Type() string {
	return "secp256k1"
}

func (a *KeyWrapper) Reset() {
	// no-op
}

func (a *KeyWrapper) String() string {
	return "<redacted>"
}

func (a *KeyWrapper) ProtoMessage() {
	// no-op
}
