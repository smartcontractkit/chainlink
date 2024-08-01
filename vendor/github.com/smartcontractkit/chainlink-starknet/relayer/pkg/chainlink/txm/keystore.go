package txm

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	starknetaccount "github.com/NethermindEth/starknet.go/account"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	adapters "github.com/smartcontractkit/chainlink-common/pkg/loop/adapters/starknet"
)

// KeystoreAdapter is a starknet-specific adaption layer to translate between the generic Loop Keystore (bytes) and
// the type specific caigo Keystore (big.Int)
type KeystoreAdapter interface {
	starknetaccount.Keystore
	Loopp() loop.Keystore
}

// keystoreAdapter implements [KeystoreAdapter].
type keystoreAdapter struct {
	looppKs loop.Keystore
}

// NewKeystoreAdapter instantiates the KeystoreAdapter interface
// The implementation requires that the given [looppKs] produces a signature [loop.Keystore.Sign]
// that is []byte representation of [adapters.Signature]
// Callers are responsible for ensuring that the given LOOPp Keystore encodes
// signatures correctly.
func NewKeystoreAdapter(lk loop.Keystore) KeystoreAdapter {
	return &keystoreAdapter{looppKs: lk}
}

var ErrBadAdapterEncoding = errors.New("failed to decode raw signature as adapter signature")

// Sign implements the caigo Keystore Sign func. Returns [ErrBadAdapterSignature] if the signature cannot be
// decoded from the [loop.Keystore] implementation
func (ca *keystoreAdapter) Sign(ctx context.Context, senderAddress string, hash *big.Int) (*big.Int, *big.Int, error) {
	raw, err := ca.looppKs.Sign(ctx, senderAddress, hash.Bytes())
	if err != nil {
		return nil, nil, fmt.Errorf("error computing loopp keystore signature: %w", err)
	}
	starknetSig, serr := adapters.SignatureFromBytes(raw)
	if serr != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrBadAdapterEncoding, serr)
	}
	return starknetSig.Ints()
}

func (ca *keystoreAdapter) Loopp() loop.Keystore {
	return ca.looppKs
}
