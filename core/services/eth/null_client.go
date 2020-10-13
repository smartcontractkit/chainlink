package eth

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// NullClient satisfies the Client but has no side effects
type NullClient struct{}

//
// Client methods
//

func (nc *NullClient) Dial(ctx context.Context) error {
	return nil
}

func (nc *NullClient) Close() {}

func (nc *NullClient) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (nc *NullClient) SendRawTx(bytes []byte) (common.Hash, error) {
	return common.Hash{}, nil
}

func (nc *NullClient) Call(result interface{}, method string, args ...interface{}) error {
	return nil
}

func (nc *NullClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return nil
}

func (nc *NullClient) HeaderByNumber(ctx context.Context, n *big.Int) (*models.Head, error) {
	return nil, nil
}

type nullSubscription struct{}

func (ns *nullSubscription) Unsubscribe()      {}
func (ns *nullSubscription) Err() <-chan error { return nil }

func (nc *NullClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return &nullSubscription{}, nil
}

func (nc *NullClient) SubscribeNewHead(ctx context.Context, ch chan<- *models.Head) (ethereum.Subscription, error) {
	return &nullSubscription{}, nil
}

//
// GethClient methods
//

func (nc *NullClient) ChainID(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}

func (nc *NullClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func (nc *NullClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return nil, nil
}

func (nc *NullClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil
}

func (nc *NullClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return nil, nil
}

func (nc *NullClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return nil, nil
}

func (nc *NullClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (nc *NullClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return nil, nil
}

func (nc *NullClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return 0, nil
}

func (nc *NullClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (nc *NullClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}

func (nc *NullClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}
