package cosmostxm

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/cosmoskey"
)

// Note we use this strictly for https://github.com/smartcontractkit/cosmos.go/blob/master/tx/txbuilder.go#L37
// i.e. inline signing txes.
var _ cryptotypes.PrivKey = KeyWrapper{}

// KeyWrapper wrapper around a cosmos transmitter key
// for use in the cosmos txbuilder and client.
type KeyWrapper struct {
	key cosmoskey.Key
}

// NewKeyWrapper create a key wrapper
func NewKeyWrapper(key cosmoskey.Key) KeyWrapper {
	return KeyWrapper{key: key}
}

// Reset nop
func (k KeyWrapper) Reset() {}

// ProtoMessage nop
func (k KeyWrapper) ProtoMessage() {}

// String nop
func (k KeyWrapper) String() string {
	return ""
}

// Bytes does not expose private key
func (k KeyWrapper) Bytes() []byte {
	return []byte{}
}

// Sign sign a message with key
func (k KeyWrapper) Sign(msg []byte) ([]byte, error) {
	return k.key.ToPrivKey().Sign(msg)
}

// PubKey get the pubkey
func (k KeyWrapper) PubKey() cryptotypes.PubKey {
	return k.key.PublicKey()
}

// Equals compare against another key
func (k KeyWrapper) Equals(a cryptotypes.LedgerPrivKey) bool {
	return k.PubKey().Address().String() == a.PubKey().Address().String()
}

// Type nop
func (k KeyWrapper) Type() string {
	return ""
}
