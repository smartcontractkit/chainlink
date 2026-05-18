package keystore

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Eth is the external interface for EthKeyStore
type Eth interface {
	CheckEnabled(ctx context.Context, address common.Address, chainID *big.Int) error
	EnabledAddressesForChain(ctx context.Context, chainID *big.Int) (addresses []common.Address, err error)
	SignTx(ctx context.Context, fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	SubscribeToKeyChanges(ctx context.Context) (ch chan struct{}, unsub func())
}
