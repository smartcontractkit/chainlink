package client

import (
	"context"
	"math/big"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

var _ Node = (*erroringNode)(nil)

type erroringNode struct {
	errMsg string
}

func (e *erroringNode) ChainID() (chainID *big.Int) { return nil }

func (e *erroringNode) Start(ctx context.Context) error { return errors.New(e.errMsg) }

func (e *erroringNode) Close() error { return nil }

func (e *erroringNode) Verify(ctx context.Context, expectedChainID *big.Int) (err error) {
	return errors.New(e.errMsg)
}

func (e *erroringNode) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) HeaderByHash(_ context.Context, _ common.Hash) (*types.Header, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) EthSubscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode) String() string {
	return "<erroring node>"
}

func (e *erroringNode) State() NodeState {
	return NodeStateUnreachable
}

func (e *erroringNode) StateAndLatest() (NodeState, int64, *utils.Big) {
	return NodeStateUnreachable, -1, nil
}

func (e *erroringNode) DeclareOutOfSync()            {}
func (e *erroringNode) DeclareInSync()               {}
func (e *erroringNode) DeclareUnreachable()          {}
func (e *erroringNode) Name() string                 { return "" }
func (e *erroringNode) NodeStates() map[int32]string { return nil }
