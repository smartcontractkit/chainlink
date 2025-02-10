package clientwrappers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-integrations/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type GethClient struct {
	*ethclient.Client
}

func NewGethClient(client *ethclient.Client) *GethClient {
	return &GethClient{
		Client: client,
	}
}

func (g *GethClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return g.Client.Client().BatchCallContext(ctx, b)
}

func (g *GethClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return g.Client.Client().CallContext(ctx, result, method, args...)
}

func (g *GethClient) CallContract(ctx context.Context, message ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := g.CallContext(ctx, &hex, "eth_call", client.ToBackwardCompatibleCallArg(message), client.ToBackwardCompatibleBlockNumArg(blockNumber))
	return hex, err
}

func (g *GethClient) HeadByNumber(ctx context.Context, number *big.Int) (*evmtypes.Head, error) {
	hexNumber := client.ToBlockNumArg(number)
	args := []interface{}{hexNumber, false}
	head := new(evmtypes.Head)
	err := g.CallContext(ctx, head, "eth_getBlockByNumber", args...)
	return head, err
}

func (g *GethClient) SendTransaction(ctx context.Context, _ *types.Transaction, attempt *types.Attempt) error {
	return g.Client.SendTransaction(ctx, attempt.SignedTransaction)
}
