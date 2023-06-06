package types

import (
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// KeyStore encompasses the subset of keystore used by txmgr
//
//go:generate mockery --quiet --name KeyStore --output ./mocks/ --case=underscore
type KeyStore[
	// Account Address type.
	ADDR types.Hashable,
	// Chain ID type
	CHAIN_ID ID,
	// Chain's sequence type. For example, EVM chains use nonce, bitcoin uses UTXO.
	SEQ Sequence,
] interface {
	CheckEnabled(address ADDR, chainID CHAIN_ID) error
	NextSequence(address ADDR, chainID CHAIN_ID, qopts ...pg.QOpt) (SEQ, error)
	EnabledAddressesForChain(chainId CHAIN_ID) ([]ADDR, error)
	IncrementNextSequence(address ADDR, chainID CHAIN_ID, currentSequence SEQ, qopts ...pg.QOpt) error
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}
