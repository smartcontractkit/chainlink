package cltest

import (
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

	"chainlink/core/assets"
	"chainlink/core/eth"
	"chainlink/core/utils"
)

// SimulatedBackendClient is an eth.SimulatedBackendClient implementation using
// a simulated blockchain backend. Note that not all RPC methods are implemented
// here.
type SimulatedBackendClient struct {
	b *backends.SimulatedBackend
	t testing.TB
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
	if err != nil {
		return nil, nil, fmt.Errorf("fourth arg to SimulatedBackendClient.Call " +
			"must be a positive *big.Int, or one of the strings \"latest\", " +
			"\"pending\", or \"earliest\"")
	}
	return &callArgs, blockNumber, nil
}

// Call mocks the ethereum client RPC calls used by chainlink.
func (c *SimulatedBackendClient) Call(result interface{}, method string,
	args ...interface{}) error {
	switch method {
	case "eth_call":
		callArgs, blockNumber, err := c.checkEthCallArgs(args)
		if err != nil {
			return err
		}
		b, err := c.b.CallContract(context.Background(), ethereum.CallMsg{
			To: &callArgs.To, Data: callArgs.Data}, blockNumber)
		if err != nil {
			return errors.Wrapf(err, "while calling contract at address %x with "+
				"data %x", callArgs.To, callArgs.Data)
		}
		switch r := result.(type) {
		case *hexutil.Bytes:
			copy(*r, b)
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

// XXX: Move these to utils.
// chainlinkEthLogFromGethLog returns a copy of l as an eth.Log. (They have
// identical fields, but the field tags differ, and the types differ slightly.)
func chainlinkEthLogFromGethLog(l types.Log) eth.Log {
	return eth.Log{Address: l.Address, Topics: l.Topics, Data: l.Data,
		BlockNumber: l.BlockNumber, TxHash: l.TxHash, TxIndex: l.TxIndex,
		BlockHash: l.BlockHash, Index: l.Index, Removed: l.Removed}
}

// GetLogs returns all logs that respect the passed filter query.
func (c *SimulatedBackendClient) GetLogs(q ethereum.FilterQuery) (logs []eth.Log,
	err error) {
	rawLogs, err := c.b.FilterLogs(context.Background(), q)
	if err != nil {
		return nil, errors.Wrapf(err, "while querying for logs with %s", q)
	}
	for _, rawLog := range rawLogs {
		logs = append(logs, chainlinkEthLogFromGethLog(rawLog))
	}
	return logs, nil
}

// SubscribeToLogs registers a subscription for push notifications of logs
// from a given address.
func (c *SimulatedBackendClient) SubscribeToLogs(ctx context.Context, channel chan<- eth.Log,
	q ethereum.FilterQuery) (eth.Subscription, error) {
	ch := make(chan types.Log)
	go func() {
		for l := range ch {
			channel <- chainlinkEthLogFromGethLog(l)
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

var balanceOfAPI abi.ABI

func init() {
	var err error
	balanceOfAPI, err = abi.JSON(strings.NewReader(balanceOfABIString))
	if err != nil {
		panic(errors.Wrapf(err, "while parsing erc20ABI"))
	}
}

// GetERC20Balance returns the balance of the given address for the token
// contract address.
func (c *SimulatedBackendClient) GetERC20Balance(address common.Address,
	contractAddress common.Address) (balance *big.Int, err error) {
	callData, err := balanceOfAPI.Pack("balanceOf", address)
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
	return balance, balanceOfAPI.Unpack(balance, "balanceOf", b)
}

// SendRawTx sends a signed transaction to the transaction pool.
func (c *SimulatedBackendClient) SendRawTx(hex string) (txHash common.Hash, err error) {
	tx, err := utils.DecodeEthereumTx(hex)
	if err != nil {
		return common.Hash{}, errors.Wrapf(err, "while sending tx %s", hex)
	}
	return tx.Hash(), c.b.SendTransaction(context.Background(), &tx)
}

// GetTxReceipt returns the transaction receipt for the given transaction hash.
func (c *SimulatedBackendClient) GetTxReceipt(receipt common.Hash) (*eth.TxReceipt, error) {
	rawReceipt, err := c.b.TransactionReceipt(context.Background(), receipt)
	if err != nil {
		return nil, errors.Wrapf(err, "while retrieving tx receipt for %s", receipt)
	}
	logs := []eth.Log{}
	for _, log := range rawReceipt.Logs {
		logs = append(logs, chainlinkEthLogFromGethLog(*log))
	}
	return &eth.TxReceipt{BlockNumber: (*utils.Big)(rawReceipt.BlockNumber),
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
			return big.NewInt(0).Add(c.currentBlockNumber(), big.NewInt(1)), nil
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
func (c *SimulatedBackendClient) GetBlockByNumber(hex string) (block eth.Block,
	err error) {
	blockNumber, err := c.blockNumber(hex)
	if err != nil {
		c.t.Fatalf("while getting block by number: %s", err)
	}
	b, err := c.b.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return eth.Block{}, errors.Wrapf(err, "while retrieving block %d",
			blockNumber)
	}
	var txs []eth.Transaction
	for _, tx := range b.Transactions() {
		txs = append(txs, eth.Transaction{
			GasPrice: hexutil.Uint64(tx.GasPrice().Uint64())})
	}
	difficulty := hexutil.Uint64(b.Difficulty().Uint64())
	return eth.Block{Transactions: txs, Difficulty: difficulty}, nil
}

// GetChainID returns the ethereum ChainID.
func (c *SimulatedBackendClient) GetChainID() (*big.Int, error) {
	// The actual chain ID is c.b.Blockchain().Config().ChainID, but here we need
	// to match the chain ID used by the testing harness.
	return big.NewInt(int64(testConfigConstants["ETH_CHAIN_ID"].(int))), nil
}

// SubscribeToNewHeads registers a subscription for push notifications of new
// blocks.
func (c *SimulatedBackendClient) SubscribeToNewHeads(ctx context.Context,
	channel chan<- eth.BlockHeader) (eth.Subscription, error) {
	ch := make(chan *types.Header)
	go func() {
		for h := range ch {
			channel <- eth.BlockHeader{ParentHash: h.ParentHash, UncleHash: h.UncleHash,
				Coinbase: h.Coinbase, Root: h.Root, TxHash: h.TxHash,
				ReceiptHash: h.ReceiptHash, Bloom: h.Bloom,
				Difficulty: hexutil.Big(*h.Difficulty), Number: hexutil.Big(*h.Number),
				GasLimit: hexutil.Uint64(h.GasLimit), GasUsed: hexutil.Uint64(h.GasUsed),
				Time:  hexutil.Big(*big.NewInt(int64(h.Time))),
				Extra: hexutil.Bytes(h.Extra), Nonce: h.Nonce, GethHash: h.Hash(),
				// ParityHash not included, because this client is strictly based on
				// go-ethereum
			}
		}
	}()
	return c.b.SubscribeNewHead(context.Background(), ch)
}
