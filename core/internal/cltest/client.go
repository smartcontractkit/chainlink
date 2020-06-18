package cltest

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// SimulatedBackendClient is an eth.SimulatedBackendClient implementation using
// a simulated blockchain backend. Note that not all RPC methods are implemented
// here.
type SimulatedBackendClient struct {
	b       *backends.SimulatedBackend
	t       testing.TB
	chainId int
}

// GethClient is a noop, solely needed to conform to GethClientWrapper interface
func (c *SimulatedBackendClient) GethClient(f func(c eth.GethClient) error) error {
	return nil
}

// Close terminates the underlying blockchain's update loop.
func (c *SimulatedBackendClient) Close() {
	c.b.Close()
}

var _ eth.Client = (*SimulatedBackendClient)(nil)

// checkEthCallArgs extracts and verifies the arguments for an eth_call RPC
func (c *SimulatedBackendClient) checkEthCallArgs(
	args []interface{}) (*eth.CallArgs, *big.Int, error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf(
			"should have two arguments after \"eth_call\", got %d", len(args))
	}
	callArgs, ok := args[0].(eth.CallArgs)
	if !ok {
		return nil, nil, fmt.Errorf("third arg to SimulatedBackendClient.Call "+
			"must be an eth.CallArgs, got %+#v", args[0])
	}
	blockNumber, err := c.blockNumber(args[1])
	if err != nil || blockNumber.Cmp(c.currentBlockNumber()) != 0 {
		return nil, nil, fmt.Errorf("fourth arg to SimulatedBackendClient.Call "+
			"must be the string \"latest\", or a *big.Int equal to current "+
			"blocknumber, got %#+v", args[1])
	}
	return &callArgs, blockNumber, nil
}

// Call mocks the ethereum client RPC calls used by chainlink, copying the
// return value into result.
func (c *SimulatedBackendClient) Call(result interface{}, method string,
	args ...interface{}) error {
	switch method {
	case "eth_call":
		callArgs, _, err := c.checkEthCallArgs(args)
		if err != nil {
			return err
		}
		callMsg := ethereum.CallMsg{To: &callArgs.To, Data: callArgs.Data}
		b, err := c.b.CallContract(context.TODO(), callMsg, nil /* always latest block */)
		if err != nil {
			return errors.Wrapf(err, "while calling contract at address %x with "+
				"data %x", callArgs.To, callArgs.Data)
		}
		switch r := result.(type) {
		case *hexutil.Bytes:
			*r = append(*r, b...)
			if !bytes.Equal(*r, b) {
				return fmt.Errorf("was passed a non-empty array, or failed to copy "+
					"answer. Expected %x = %x", *r, b)
			}
			return nil
		default:
			return fmt.Errorf("first arg to SimulatedBackendClient.Call is an "+
				"unrecognized type: %T; add processing logic for it here", result)
		}
	default:
		return fmt.Errorf("second arg to SimulatedBackendClient.Call is an RPC "+
			"API method which has not yet been implemented: %s. Add processing for "+
			"it here", method)
	}
}

// Subscribe is a dummy method present only to satisfy the eth.Client interface.
// The original method is for subscribing to events observed by the RPC client,
// but for a simulated backend that makes no sense.
func (c *SimulatedBackendClient) Subscribe(ctx context.Context, namespace interface{},
	channelAndArgs ...interface{}) (eth.Subscription, error) {
	// if these are needed, there are subscribe methods on SimulatedBackend
	panic("unimplemented")
}

// GetLogs returns all logs that respect the passed filter query.
func (c *SimulatedBackendClient) GetLogs(q ethereum.FilterQuery) (logs []models.Log,
	err error) {
	rawLogs, err := c.b.FilterLogs(context.Background(), q)
	if err != nil {
		return nil, errors.Wrapf(err, "while querying for logs with %s", q)
	}
	for _, rawLog := range rawLogs {
		logs = append(logs, ChainlinkEthLogFromGethLog(rawLog))
	}
	return logs, nil
}

// SubscribeToLogs registers a subscription for push notifications of logs
// from a given address.
func (c *SimulatedBackendClient) SubscribeToLogs(ctx context.Context, channel chan<- models.Log,
	q ethereum.FilterQuery) (eth.Subscription, error) {
	ch := make(chan types.Log)
	go func() {
		for l := range ch {
			channel <- ChainlinkEthLogFromGethLog(l)
		}
	}()
	return c.b.SubscribeFilterLogs(ctx, q, ch)
}

// currentBlockNumber returns index of *pending* block in simulated blockchain
func (c *SimulatedBackendClient) currentBlockNumber() *big.Int {
	return c.b.Blockchain().CurrentBlock().Number()
}

// GetNonce returns the nonce (transaction count) for a given address.
func (c *SimulatedBackendClient) GetNonce(address common.Address) (nonce uint64, err error) {
	return c.b.NonceAt(context.Background(), address, c.currentBlockNumber())
}

// GetEthBalance returns the balance of the given addresses in Ether.
func (c *SimulatedBackendClient) GetEthBalance(address common.Address,
) (balance *assets.Eth, err error) {
	b, err := c.b.BalanceAt(context.Background(), address, c.currentBlockNumber())
	ab := assets.Eth(*b)
	return &ab, err
}

var balanceOfABIString string = `[
  {
    "constant": true,
    "inputs": [
      {
        "name": "_owner",
        "type": "address"
      }
    ],
    "name": "balanceOf",
    "outputs": [
      {
        "name": "balance",
        "type": "uint256"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  }
]`

var balanceOfABI abi.ABI

func init() {
	var err error
	balanceOfABI, err = abi.JSON(strings.NewReader(balanceOfABIString))
	if err != nil {
		panic(errors.Wrapf(err, "while parsing erc20ABI"))
	}
}

// GetERC20Balance returns the balance of the given address for the token
// contract address.
func (c *SimulatedBackendClient) GetERC20Balance(address common.Address,
	contractAddress common.Address) (balance *big.Int, err error) {
	callData, err := balanceOfABI.Pack("balanceOf", address)
	if err != nil {
		return nil, errors.Wrapf(err, "while seeking the ERC20 balance of %s on %s",
			address, contractAddress)
	}
	b, err := c.b.CallContract(context.Background(), ethereum.CallMsg{
		To: &contractAddress, Data: callData},
		c.currentBlockNumber())
	if err != nil {
		return nil, errors.Wrapf(err, "while calling ERC20 balanceOf method on %s "+
			"for balance of %s", contractAddress, address)
	}
	balance = new(big.Int)
	return balance, balanceOfABI.Unpack(balance, "balanceOf", b)
}

// SendRawTx sends a signed transaction to the transaction pool.
func (c *SimulatedBackendClient) SendRawTx(
	txBytes []byte) (txHash common.Hash, err error) {
	tx, err := utils.DecodeEthereumTx(hexutil.Encode(txBytes))
	if err != nil {
		logger.Errorf("could not deserialize transaction: %x", txBytes)
		return common.Hash{}, errors.Wrapf(err, "while sending tx %x", txBytes)
	}
	if err = c.b.SendTransaction(context.Background(), &tx); err == nil {
		c.b.Commit()
	}
	return tx.Hash(), err
}

// GetTxReceipt returns the transaction receipt for the given transaction hash.
func (c *SimulatedBackendClient) GetTxReceipt(
	receipt common.Hash) (*models.TxReceipt, error) {
	rawReceipt, err := c.b.TransactionReceipt(context.Background(), receipt)
	if err != nil {
		return nil, errors.Wrapf(err, "while retrieving tx receipt for %s", receipt)
	}
	if rawReceipt == nil {
		// Calling code depends on getting empty TxReceipt, rather than nil
		return &models.TxReceipt{}, nil
	}
	logs := []models.Log{}
	for _, log := range rawReceipt.Logs {
		logs = append(logs, ChainlinkEthLogFromGethLog(*log))
	}
	return &models.TxReceipt{BlockNumber: (*utils.Big)(rawReceipt.BlockNumber),
		BlockHash: &rawReceipt.BlockHash, Hash: receipt, Logs: logs}, nil
}

// GetBlockHeight returns height of latest block in the simulated blockchain.
func (c *SimulatedBackendClient) GetBlockHeight() (height uint64, err error) {
	return c.currentBlockNumber().Uint64() - 1, nil
}

func (c *SimulatedBackendClient) blockNumber(
	number interface{}) (blockNumber *big.Int, err error) {
	switch n := number.(type) {
	case string:
		switch n {
		case "latest":
			return c.currentBlockNumber(), nil
		case "earliest":
			return big.NewInt(0), nil
		case "pending":
			panic("not implemented") // I don't understand the semantics of this.
			// return big.NewInt(0).Add(c.currentBlockNumber(), big.NewInt(1)), nil
		default:
			blockNumber, err = utils.HexToUint256(n)
			if err != nil {
				return nil, errors.Wrapf(err, "while parsing '%s' as hex-encoded"+
					"block number", n)
			}
			return blockNumber, nil
		}
	case *big.Int:
		if n.Sign() < 0 {
			return nil, fmt.Errorf("block number musts be non-negative")
		}
		return n, nil
	}
	panic("can never reach here")
}

// GetBlockByNumber returns the block for the passed hex, or "latest",
// "earliest", "pending". Includes all transactions
func (c *SimulatedBackendClient) GetBlockByNumber(hex string) (block models.Block,
	err error) {
	blockNumber, err := c.blockNumber(hex)
	if err != nil {
		c.t.Fatalf("while getting block by number: %s", err)
	}
	b, err := c.b.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return models.Block{}, errors.Wrapf(err, "while retrieving block %d",
			blockNumber)
	}
	var txs []models.Transaction
	for _, tx := range b.Transactions() {
		txs = append(txs, models.Transaction{
			GasPrice: hexutil.Uint64(tx.GasPrice().Uint64())})
	}
	return models.Block{Number: hexutil.Uint64(blockNumber.Uint64()),
		Transactions: txs}, nil
}

// GetChainID returns the ethereum ChainID.
func (c *SimulatedBackendClient) GetChainID() (*big.Int, error) {
	// The actual chain ID is c.b.Blockchain().Config().ChainID, but here we need
	// to match the chain ID used by the testing harness.
	return big.NewInt(int64(c.chainId)), nil
}

// SubscribeToNewHeads registers a subscription for push notifications of new
// blocks.
func (c *SimulatedBackendClient) SubscribeToNewHeads(ctx context.Context,
	channel chan<- types.Header) (eth.Subscription, error) {
	ch := make(chan *types.Header)
	go func() {
		for h := range ch {
			channel <- *h
		}
		close(channel)
	}()
	return c.b.SubscribeNewHead(context.Background(), ch)
}

// GetLatestBlock returns the last committed block of the best blockchain the
// blockchain node is aware of.
func (c *SimulatedBackendClient) GetLatestBlock() (models.Block, error) {
	height, err := c.GetBlockHeight()
	if err != nil {
		return models.Block{}, errors.Wrap(err, "while getting latest block")
	}
	return c.GetBlockByNumber(common.BigToHash(big.NewInt(int64(height))).Hex())
}
