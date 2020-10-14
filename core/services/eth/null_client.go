package eth

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// NullClient satisfies the Client but has no side effects
type NullClient struct{}

//
// Client methods
//

func (nc *NullClient) Dial(ctx context.Context) error {
	logger.Debug("NullClient#Dial")
	return nil
}

func (nc *NullClient) Close() {
	logger.Debug("NullClient#Close")
}

func (nc *NullClient) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	logger.Debug("NullClient#GetERC20Balance")
	return big.NewInt(0), nil
}

func (nc *NullClient) SendRawTx(bytes []byte) (common.Hash, error) {
	logger.Debug("NullClient#SendRawTx")
	return common.Hash{}, nil
}

func (nc *NullClient) Call(result interface{}, method string, args ...interface{}) error {
	logger.Debug("NullClient#Call")
	return nil
}

func (nc *NullClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	logger.Debug("NullClient#CallContext")
	return nil
}

func (nc *NullClient) HeaderByNumber(ctx context.Context, n *big.Int) (*models.Head, error) {
	logger.Debug("NullClient#HeaderByNumber")
	return nil, nil
}

type nullSubscription struct{}

func (ns *nullSubscription) Unsubscribe() {
	logger.Debug("NullClient nullSubscription#Unsubscribe")
}

func (ns *nullSubscription) Err() <-chan error {
	logger.Debug("NullClient nullSubscription#Err")
	return nil
}

func (nc *NullClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	logger.Debug("NullClient#SubscribeFilterLogs")
	return &nullSubscription{}, nil
}

func (nc *NullClient) SubscribeNewHead(ctx context.Context, ch chan<- *models.Head) (ethereum.Subscription, error) {
	logger.Debug("NullClient#SubscribeNewHead")
	return &nullSubscription{}, nil
}

//
// GethClient methods
//

func (nc *NullClient) ChainID(ctx context.Context) (*big.Int, error) {
	logger.Debug("NullClient#ChainID")
	return big.NewInt(1), nil
}

func (nc *NullClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	logger.Debug("NullClient#SendTransaction")
	return nil
}

func (nc *NullClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	logger.Debug("NullClient#PendingCodeAt")
	return nil, nil
}

func (nc *NullClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	logger.Debug("NullClient#PendingNonceAt")
	return 0, nil
}

func (nc *NullClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	logger.Debug("NullClient#TransactionReceipt")
	return nil, nil
}

func (nc *NullClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	logger.Debug("NullClient#BlockByNumber")
	return nil, nil
}

func (nc *NullClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	logger.Debug("NullClient#BalanceAt")
	return big.NewInt(0), nil
}

func (nc *NullClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	logger.Debug("NullClient#FilterLogs")
	return nil, nil
}

func (nc *NullClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	logger.Debug("NullClient#EstimateGas")
	return 0, nil
}

func (nc *NullClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	logger.Debug("NullClient#SuggestGasPrice")
	return big.NewInt(0), nil
}

func (nc *NullClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	logger.Debug("NullClient#CallContract")
	return nil, nil
}

func (nc *NullClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	logger.Debug("NullClient#CodeAt")
	return nil, nil
}
