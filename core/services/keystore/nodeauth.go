package keystore

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

func BuildNodeAuth(
	ctx context.Context,
	keyStoreCSA CSA,
	keyStoreP2P P2P,
) (
	signer *core.Ed25519Signer,
	csaPubkey ed25519.PublicKey,
	p2pID types.PeerID,
	err error,
) {
	csaKey, err := GetDefault(ctx, keyStoreCSA)
	if err != nil {
		return nil, nil, types.PeerID{}, err
	}

	p2pKey, err := GetDefault(ctx, keyStoreP2P)
	if err != nil {
		return nil, nil, types.PeerID{}, err
	}

	// Create ed25519 signer from the node's csa private key
	signFn := func(ctx context.Context, account string, data []byte) (signed []byte, err error) {
		return csaKey.Sign(rand.Reader, data, crypto.Hash(0))
	}

	signer, err = core.NewEd25519Signer(hex.EncodeToString(csaKey.PublicKey), signFn)
	if err != nil {
		return nil, nil, types.PeerID{}, err
	}
	csaPubkey = csaKey.PublicKey
	p2pID = types.PeerID(p2pKey.PeerID())
	return
}
