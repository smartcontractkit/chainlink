package types

import (
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// KeyStore encompasses the subset of keystore used by txmgr
//
//go:generate mockery --quiet --name KeyStore --output ./mocks/ --case=underscore
type KeyStore[ADDR types.Hashable, ID any, TX any, META any] interface {
	CheckEnabled(address ADDR, chainID ID) error
	NextSequence(address ADDR, chainID ID, qopts ...pg.QOpt) (META, error)
	EnabledAddressesForChain(chainId ID) ([]ADDR, error)
	IncrementNextSequence(address ADDR, chainID ID, currentNonce META, qopts ...pg.QOpt) error
	SignTx(fromAddress ADDR, tx *TX, chainID ID) (*TX, error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}
