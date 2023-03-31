package types

import (
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// KeyStore encompasses the subset of keystore used by txmgr
//
//go:generate mockery --quiet --name KeyStore --output ./mocks/ --case=underscore
type KeyStore[ADDR types.Hashable, ID any, S any] interface {
	CheckEnabled(address ADDR, chainID ID) error
	NextSequence(address ADDR, chainID ID, qopts ...pg.QOpt) (S, error)
	EnabledAddressesForChain(chainId ID) ([]ADDR, error)
	IncrementNextSequence(address ADDR, chainID ID, currentNonce S, qopts ...pg.QOpt) error
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())
}
