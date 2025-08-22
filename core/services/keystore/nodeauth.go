package keystore

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

func BuildNodeAuth(
	ctx context.Context,
	keyStoreCSA CSA,
) (
	*core.Ed25519Signer,
	ed25519.PublicKey,
	error,
) {
	csaKey, err := GetDefault(ctx, keyStoreCSA)
	if err != nil {
		return nil, nil, err
	}

	// Create ed25519 signer from the node's csa private key
	signFn := func(ctx context.Context, account string, data []byte) (signed []byte, err error) {
		return csaKey.Sign(rand.Reader, data, crypto.Hash(0))
	}

	signer, err := core.NewEd25519Signer(hex.EncodeToString(csaKey.PublicKey), signFn)
	if err != nil {
		return nil, nil, err
	}

	return signer, csaKey.PublicKey, nil
}
