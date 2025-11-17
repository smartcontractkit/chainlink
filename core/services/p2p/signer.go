package p2p

import (
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type signer struct {
	keystoreP2P keystore.P2P
	peerID      p2pkey.PeerID
	privateKey  crypto.Signer
}

var _ types.Signer = &signer{}

// We can't access the Keystore until it is started and unlocked. The owner of the
// Signer object needs to call Initialize() when Keystore is ready and before calling Sign().
func NewSigner(keystoreP2P keystore.P2P, peerID p2pkey.PeerID) *signer {
	return &signer{
		keystoreP2P: keystoreP2P,
		peerID:      peerID,
	}
}

func (s *signer) Initialize() error {
	key, err := s.keystoreP2P.GetOrFirst(s.peerID)
	if err != nil {
		return fmt.Errorf("failed to get P2P key from keystore: %w", err)
	}
	s.privateKey = key
	return nil
}

func (s *signer) Sign(msg []byte) ([]byte, error) {
	if s.privateKey == nil {
		return nil, errors.New("private key not set")
	}
	return s.privateKey.Sign(rand.Reader, msg, crypto.Hash(0))
}
