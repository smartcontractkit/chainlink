package terratxm

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
)

type PrivKey struct {
	key terrakey.Key
}

func NewPrivKey(key terrakey.Key) PrivKey {
	return PrivKey{key: key}
}

// protobuf methods (don't do anything)
func (k PrivKey) Reset()        {}
func (k PrivKey) ProtoMessage() {}
func (k PrivKey) String() string {
	return ""
}

func (k PrivKey) Bytes() []byte {
	return []byte{} // does not expose private key
}
func (k PrivKey) Sign(msg []byte) ([]byte, error) {
	return k.key.Sign(msg)
}
func (k PrivKey) PubKey() cryptotypes.PubKey {
	return k.key.PublicKey()
}
func (k PrivKey) Equals(a cryptotypes.LedgerPrivKey) bool {
	return k.PubKey().Address().String() == a.PubKey().Address().String()
}
func (k PrivKey) Type() string {
	return ""
}
