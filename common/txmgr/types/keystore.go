package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-framework/chains"
)

// KeyStore encompasses the subset of keystore used by txmgr
type KeyStore[
	// Account Address type.
	ADDR chains.Hashable,
	// Chain ID type
	CHAIN_ID chains.ID,
	// Chain's sequence type. For example, EVM chains use nonce, bitcoin uses UTXO.
	SEQ chains.Sequence,
] interface {
	CheckEnabled(ctx context.Context, address ADDR, chainID CHAIN_ID) error
	EnabledAddressesForChain(ctx context.Context, chainId CHAIN_ID) ([]ADDR, error)
	SubscribeToKeyChanges(ctx context.Context) (ch chan struct{}, unsub func())
}
