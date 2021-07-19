package eth

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// NullClient satisfies the Client but has no side effects
type NullClient struct{
	l logger.Logger
}

// NullClientChainID the ChainID that nullclient will return
// 0 is never used as a real chain ID so makes sense as a dummy value here
const NullClientChainID = 0

//
// Client methods
//

func (nc *NullClient) Dial(ctx context.Context) error {
	nc.l.Debug("NullClient#Dial")
	return nil
}

func (nc *NullClient) Close() {
	nc.l.Debug("NullClient#Close")
}

func (nc *NullClient) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	nc.l.Debug("NullClient#GetERC20Balance")
	return big.NewInt(0), nil
}

func (nc *NullClient) GetLINKBalance(linkAddress common.Address, address common.Address) (*assets.Link, error) {
	nc.l.Debug("NullClient#GetLINKBalance")
	return assets.NewLink(0), nil
}

func (nc *NullClient) GetEthBalance(context.Context, common.Address, *big.Int) (*assets.Eth, error) {
	nc.l.Debug("NullClient#GetEthBalance")
	return assets.NewEth(0), nil
}

func (nc *NullClient) Call(result interface{}, method string, args ...interface{}) error {
	nc.l.Debug("NullClient#Call")
	return nil
}

func (nc *NullClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	nc.l.Debug("NullClient#CallContext")
	return nil
}

func (nc *NullClient) HeadByNumber(ctx context.Context, n *big.Int) (*models.Head, error) {
	nc.l.Debug("NullClient#HeadByNumber")
	return &models.Head{}, nil
}

type nullSubscription struct{
	l logger.Logger
}

func (ns *nullSubscription) Unsubscribe() {
	ns.l.Debug("NullClient nullSubscription#Unsubscribe")
}

func (ns *nullSubscription) Err() <-chan error {
	ns.l.Debug("NullClient nullSubscription#Err")
	return nil
}

func (nc *NullClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	nc.l.Debug("NullClient#SubscribeFilterLogs")
	return &nullSubscription{l: nc.l}, nil
}

func (nc *NullClient) SubscribeNewHead(ctx context.Context, ch chan<- *models.Head) (ethereum.Subscription, error) {
	nc.l.Debug("NullClient#SubscribeNewHead")
	return &nullSubscription{l: nc.l}, nil
}

//
// GethClient methods
//

func (nc *NullClient) ChainID(ctx context.Context) (*big.Int, error) {
	nc.l.Debug("NullClient#ChainID")
	return big.NewInt(NullClientChainID), nil
}

func (nc *NullClient) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	nc.l.Debug("NullClient#HeaderByNumber")
	return nil, nil
}

func (nc *NullClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	nc.l.Debug("NullClient#SendTransaction")
	return nil
}

func (nc *NullClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	nc.l.Debug("NullClient#PendingCodeAt")
	return nil, nil
}

func (nc *NullClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	nc.l.Debug("NullClient#PendingNonceAt")
	return 0, nil
}

func (nc *NullClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	nc.l.Debug("NullClient#NonceAt")
	return 0, nil
}

func (nc *NullClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	nc.l.Debug("NullClient#TransactionReceipt")
	return nil, nil
}

func (nc *NullClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	nc.l.Debug("NullClient#BlockByNumber")
	return nil, nil
}

func (nc *NullClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	nc.l.Debug("NullClient#BalanceAt")
	return big.NewInt(0), nil
}

func (nc *NullClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	nc.l.Debug("NullClient#FilterLogs")
	return nil, nil
}

func (nc *NullClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	nc.l.Debug("NullClient#EstimateGas")
	return 0, nil
}

func (nc *NullClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	nc.l.Debug("NullClient#SuggestGasPrice")
	return big.NewInt(0), nil
}

func (nc *NullClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	nc.l.Debug("NullClient#CallContract")
	return nil, nil
}

func (nc *NullClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	nc.l.Debug("NullClient#CodeAt")
	return nil, nil
}

func (nc *NullClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return nil
}

func (nc *NullClient) RoundRobinBatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return nil
}

func (nc *NullClient) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	return nil, nil
}
