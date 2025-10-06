package crypto

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type P2PKey struct {
	EncryptedJSON []byte
	PeerID        p2pkey.PeerID
	Password      string
}

func NewP2PKey(password string) (*P2PKey, error) {
	key, err := p2pkey.NewV2()
	if err != nil {
		return nil, err
	}
	d, err := key.ToEncryptedJSON(password, utils.DefaultScryptParams)
	if err != nil {
		return nil, err
	}

	return &P2PKey{
		EncryptedJSON: d,
		PeerID:        key.PeerID(),
		Password:      password,
	}, nil
}
