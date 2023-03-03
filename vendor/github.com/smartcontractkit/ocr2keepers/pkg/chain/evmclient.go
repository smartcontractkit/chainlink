package chain

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

// evmClient expends the base EVM client by splitting batch calls into sub-batches
type evmClient struct {
	*ethclient.Client
	chHead    chan types.BlockKey
	rpcClient *rpc.Client
	batchSize int
}

// NewEVMClient is the constructor of evmClient
func NewEVMClient(client *rpc.Client, batchSize int) types.EVMClient {
	c := &evmClient{
		Client:    ethclient.NewClient(client),
		chHead:    make(chan types.BlockKey, 1000),
		rpcClient: client,
		batchSize: batchSize,
	}

	c.subscribe(context.Background())

	return c
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *evmClient) HeaderByNumber(ctx context.Context, number *big.Int) (*ethtypes.Header, error) {
	var head *ethtypes.Header
	err := ec.rpcClient.CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}

func (ec *evmClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	batches := ec.createBatches(b)

	var wg sync.WaitGroup
	errs := make([]error, len(batches))
	for i := range batches {
		wg.Add(1)
		go func(idx int, batch []rpc.BatchElem) {
			errs[idx] = ec.rpcClient.BatchCallContext(ctx, batch)
			wg.Done()
		}(i, batches[i])
	}

	wg.Wait()

	return errors.Wrap(multierr.Combine(errs...), "batch call error")
}

func (ec *evmClient) HeadTicker() chan types.BlockKey {
	return ec.chHead
}

func (ec *evmClient) subscribe(ctx context.Context) {
	ch := make(chan *ethtypes.Header, 1)
	sub, err := ec.SubscribeNewHead(ctx, ch)
	if err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-sub.Err():
			if sub, err = ec.SubscribeNewHead(ctx, ch); err != nil {
				return
			}
		case head := <-ch:
			ec.chHead <- BlockKey(head.Number.String())
		}
	}
}

func (ec *evmClient) createBatches(b []rpc.BatchElem) (batches [][]rpc.BatchElem) {
	for i := 0; i < len(b); i += ec.batchSize {
		j := i + ec.batchSize
		if j > len(b) {
			j = len(b)
		}
		batches = append(batches, b[i:j])
	}
	return
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
