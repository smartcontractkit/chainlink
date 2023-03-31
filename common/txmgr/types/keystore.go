package types

import "github.com/smartcontractkit/chainlink/v2/core/services/pg"

// KeyStore encompasses the subset of keystore used by txmgr
type KeyStore[ADDR any, ID any, S any] interface {
	CheckEnabled(address ADDR, chainID ID) error
	NextSequence(address ADDR, chainID ID, qopts ...pg.QOpt) (S, error)
	EnabledAddressesForChain(chainId ID) ([]ADDR, error)
	IncrementNextSequence(address ADDR, chainID ID, currentSequence S, qopts ...pg.QOpt) error
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}
