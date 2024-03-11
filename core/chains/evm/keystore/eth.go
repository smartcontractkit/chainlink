package keystore

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// Eth is the external interface for EthKeyStore
//
//go:generate mockery --quiet --name Eth --output mocks/ --case=underscore
type Eth interface {
	Add(ctx context.Context, address common.Address, chainID *big.Int, qopts ...pg.QOpt) error
	CheckEnabled(ctx context.Context, address common.Address, chainID *big.Int) error
	EnabledAddressesForChain(ctx context.Context, chainID *big.Int) (addresses []common.Address, err error)
	SignTx(ctx context.Context, fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	SubscribeToKeyChanges(ctx context.Context) (ch chan struct{}, unsub func())
}
