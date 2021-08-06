package eth

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

//go:generate mockery --name Client --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name Client --output mocks/ --case=underscore
//go:generate mockery --name Subscription --output ../../internal/mocks/ --case=underscore

// Client is the interface used to interact with an ethereum node.
type Client interface {
	Dial(ctx context.Context) error
	Close()

	GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error)
	GetLINKBalance(linkAddress common.Address, address common.Address) (*assets.Link, error)
	GetEthBalance(ctx context.Context, account common.Address, blockNumber *big.Int) (*assets.Eth, error)

	// Wrapped RPC methods
	Call(result interface{}, method string, args ...interface{}) error
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	RoundRobinBatchCallContext(ctx context.Context, b []rpc.BatchElem) error

	// HeadByNumber is a reimplemented version of HeaderByNumber due to a
	// difference in how block header hashes are calculated by Parity nodes
	// running on Kovan.  We have to return our own wrapper type to capture the
	// correct hash from the RPC response.
	HeadByNumber(ctx context.Context, n *big.Int) (*models.Head, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *models.Head) (ethereum.Subscription, error)

	// Wrapped Geth client methods
	ChainID(ctx context.Context) (*big.Int, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)

	// bind.ContractBackend methods
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

// This interface only exists so that we can generate a mock for it.  It is
// identical to `ethereum.Subscription`.
type Subscription interface {
	Err() <-chan error
	Unsubscribe()
}

// DefaultQueryCtx returns a context with a sensible sanity limit timeout for
// queries to the eth node
func DefaultQueryCtx(ctxs ...context.Context) (ctx context.Context, cancel context.CancelFunc) {
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, 15*time.Second)
}

// client represents an abstract client that manages connections to
// multiple ethereum nodes
type client struct {
	primary     *node
	secondaries []*secondarynode
	mocked      bool

	roundRobinCount uint32
}

var _ Client = (*client)(nil)

func NewClient(rpcUrl string, rpcHTTPURL *url.URL, secondaryRPCURLs []url.URL) (*client, error) {
	parsed, err := url.ParseRequestURI(rpcUrl)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "ws" && parsed.Scheme != "wss" {
		return nil, errors.Errorf("ethereum url scheme must be websocket: %s", parsed.String())
	}

	c := client{}

	// for now only one primary is supported
	c.primary = newNode(*parsed, rpcHTTPURL, "eth-primary-0")

	for i, url := range secondaryRPCURLs {
		if url.Scheme != "http" && url.Scheme != "https" {
			return nil, errors.Errorf("secondary ethereum rpc url scheme must be http(s): %s", url.String())
		}
		s := newSecondaryNode(url, fmt.Sprintf("eth-secondary-%d", i))
		c.secondaries = append(c.secondaries, s)
	}
	return &c, nil
}

func (client *client) Dial(ctx context.Context) error {
	if client.mocked {
		return nil
	}
	if err := client.primary.Dial(ctx); err != nil {
		return errors.Wrap(err, "Failed to dial primary client")
	}

	for _, s := range client.secondaries {
		err := s.Dial()
		if err != nil {
			return errors.Wrapf(err, "Failed to dial secondary client: %v", s.uri)
		}
	}
	return nil
}

func (client *client) Close() {
	client.primary.Close()
}

// CallArgs represents the data used to call the balance method of a contract.
// "To" is the address of the ERC contract. "Data" is the message sent
// to the contract.
type CallArgs struct {
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

// GetERC20Balance returns the balance of the given address for the token contract address.
func (client *client) GetERC20Balance(address common.Address, contractAddress common.Address) (*big.Int, error) {
	result := ""
	numLinkBigInt := new(big.Int)
	functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
	data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(address.Bytes(), utils.EVMWordByteLen))
	args := CallArgs{
		To:   contractAddress,
		Data: data,
	}
	err := client.Call(&result, "eth_call", args, "latest")
	if err != nil {
		return numLinkBigInt, err
	}
	numLinkBigInt.SetString(result, 0)
	return numLinkBigInt, nil
}

// GetLINKBalance returns the balance of LINK at the given address
func (client *client) GetLINKBalance(linkAddress common.Address, address common.Address) (*assets.Link, error) {
	balance, err := client.GetERC20Balance(address, linkAddress)
	if err != nil {
		return assets.NewLink(0), err
	}
	return (*assets.Link)(balance), nil
}

func (client *client) GetEthBalance(ctx context.Context, account common.Address, blockNumber *big.Int) (*assets.Eth, error) {
	balance, err := client.BalanceAt(ctx, account, blockNumber)
	if err != nil {
		return assets.NewEth(0), err
	}
	return (*assets.Eth)(balance), nil
}

// We wrap the GethClient's `TransactionReceipt` method so that we can ignore the error that arises
// when we're talking to a Parity node that has no receipt yet.
func (client *client) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	receipt, err = client.primary.TransactionReceipt(ctx, txHash)

	if err != nil && strings.Contains(err.Error(), "missing required field") {
		return nil, ethereum.NotFound
	}
	return
}

func (client *client) ChainID(ctx context.Context) (*big.Int, error) {
	return client.primary.ChainID(ctx)
}

func (client *client) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	return client.primary.HeaderByNumber(ctx, n)
}

// SendTransaction also uses the secondary HTTP RPC URLs if set
func (client *client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, s := range client.secondaries {
		// Parallel send to secondary node
		wg.Add(1)
		go func(s *secondarynode) {
			defer wg.Done()
			err := NewSendError(s.SendTransaction(ctx, tx))
			if err == nil || err.IsNonceTooLowError() || err.IsTransactionAlreadyInMempool() {
				// Nonce too low or transaction known errors are expected since
				// the primary SendTransaction may well have succeeded already
				return
			}
			logger.Warnw("secondary eth client returned error", "err", err, "tx", tx)
		}(s)
	}

	return client.primary.SendTransaction(ctx, tx)
}

func (client *client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return client.primary.PendingNonceAt(ctx, account)
}

func (client *client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return client.primary.NonceAt(ctx, account, blockNumber)
}

func (client *client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return client.primary.PendingCodeAt(ctx, account)
}

func (client *client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return client.primary.EstimateGas(ctx, call)
}

func (client *client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return client.primary.SuggestGasPrice(ctx)
}

func (client *client) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return client.primary.CallContract(ctx, msg, blockNumber)
}

func (client *client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return client.primary.CodeAt(ctx, account, blockNumber)
}

func (client *client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return client.primary.BlockByNumber(ctx, number)
}

func (client *client) HeadByNumber(ctx context.Context, number *big.Int) (head *models.Head, err error) {
	hex := toBlockNumArg(number)
	err = client.primary.CallContext(ctx, &head, "eth_getBlockByNumber", hex, false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

func (client *client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return client.primary.BalanceAt(ctx, account, blockNumber)
}

func (client *client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return client.primary.FilterLogs(ctx, q)
}

func (client *client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	logger.Debugw("eth.Client#SubscribeFilterLogs(...)",
		"q", q,
	)
	return client.primary.SubscribeFilterLogs(ctx, q, ch)
}

func (client *client) SubscribeNewHead(ctx context.Context, ch chan<- *models.Head) (ethereum.Subscription, error) {
	return client.primary.EthSubscribe(ctx, ch, "newHeads")
}

func (client *client) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (ethereum.Subscription, error) {
	return client.primary.EthSubscribe(ctx, channel, args...)
}

func (client *client) Call(result interface{}, method string, args ...interface{}) error {
	ctx, cancel := DefaultQueryCtx()
	defer cancel()
	return client.primary.CallContext(ctx, result, method, args...)
}

func (client *client) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return client.primary.CallContext(ctx, result, method, args...)
}

func (client *client) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return client.primary.BatchCallContext(ctx, b)
}

// RoundRobinBatchCallContext rotates through Primary and all Secondaries, changing node on each call
func (client *client) RoundRobinBatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	nSecondaries := len(client.secondaries)
	if nSecondaries == 0 {
		return client.BatchCallContext(ctx, b)
	}

	// NOTE: AddUint32 returns the number after addition, so we must -1 to get the "current" count
	count := atomic.AddUint32(&client.roundRobinCount, 1) - 1
	// idx 0 indicates the primary, subsequent indices represent secondaries
	rr := int(count % uint32(nSecondaries+1))

	if rr == 0 {
		return client.BatchCallContext(ctx, b)
	}
	return client.secondaries[rr-1].BatchCallContext(ctx, b)
}

func (client *client) SuggestGasTipCap(ctx context.Context) (tipCap *big.Int, err error) {
	return client.primary.SuggestGasTipCap(ctx)
}
