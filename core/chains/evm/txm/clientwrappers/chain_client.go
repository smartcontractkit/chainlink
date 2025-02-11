package clientwrappers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-integrations/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type ChainClient struct {
	c client.Client
}

func NewChainClient(client client.Client) *ChainClient {
	return &ChainClient{c: client}
}

func (c *ChainClient) NonceAt(ctx context.Context, address common.Address, blockNumber *big.Int) (uint64, error) {
	return c.c.NonceAt(ctx, address, blockNumber)
}

func (c *ChainClient) PendingNonceAt(ctx context.Context, address common.Address) (uint64, error) {
	return c.c.PendingNonceAt(ctx, address)
}

func (c *ChainClient) SendTransaction(ctx context.Context, _ *types.Transaction, attempt *types.Attempt) error {
	return c.c.SendTransaction(ctx, attempt.SignedTransaction)
}
