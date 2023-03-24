package types

import "github.com/smartcontractkit/chainlink/core/services/pg"

// KeyStore encompasses the subset of keystore used by txmgr
type KeyStore[ADDR any, ID any, TX any, META any] interface {
	CheckEnabled(address ADDR, chainID ID) error
	GetNextMetadata(address ADDR, chainID ID, qopts ...pg.QOpt) (META, error)
	GetEnabledAddressesForChain(chainId ID) ([]ADDR, error)
	IncrementNextMetadata(address ADDR, chainID ID, currentNonce META, qopts ...pg.QOpt) error
	SignTx(fromAddress ADDR, tx *TX, chainID ID) (*TX, error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}
