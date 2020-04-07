// package simclient is an eth.Client implementation using a simulated blockchain
// backend.
//
// This client is incompatible with contracts which use the
// core/services/eth.ConnectedContract interface, as that makes actual RPC
// calls, which are not supported here. Fixing this would be a matter of
// implementing the Client.Call method here to deal with "eth_call" calls.
package cltest

import (
	"context"
	"math/big"
	"strings"

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

// SimulatedBackendClient is an eth.SimulatedBackendClient implementation using a simulated blockchain backend.
type SimulatedBackendClient struct{ b *backends.SimulatedBackend }

// Close terminates the underlying blockchain's update loop.
func (c *SimulatedBackendClient) Close() {
	c.b.Close()
}

var _ eth.Client = (*SimulatedBackendClient)(nil)

// Call is a dummy method present only to satisfy the eth.Client interface. The
// original method is for sending an RPC call over the client, but for a
// simulated backend that makes no sense.
func (c *SimulatedBackendClient) Call(result interface{}, method string, args ...interface{},
) error {
	panic(
		"unimplemented; there is no actual RPC mechanism on a simulated blockchain")
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
		From: common.Address{}, To: &contractAddress, Data: callData},
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

// GetBlockByNumber returns the block for the passed hex, or "latest",
// "earliest", "pending". Includes all transactions
func (c *SimulatedBackendClient) GetBlockByNumber(hex string) (block eth.Block, err error) {
	var blockNumber *big.Int
	switch hex {
	case "latest":
		height, err := c.GetBlockHeight()
		if err != nil {
			return eth.Block{}, errors.Wrapf(err, "while getting latest block from "+
				"simulated blockchain")
		}
		blockNumber = big.NewInt(int64(height))
	case "earliest":
		blockNumber = big.NewInt(0)
	case "pending":
		blockNumber = c.currentBlockNumber()
	default:
		blockNumber, err = utils.HexToUint256(hex)
		if err != nil {
			return eth.Block{}, errors.Wrapf(err, "while parsing '%s' as hex-encoded"+
				"block number", hex)
		}
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
	return c.b.Blockchain().Config().ChainID, nil
}

// SubscribeToNewHeads registers a subscription for push notifications of new
// blocks.
func (c *SimulatedBackendClient) SubscribeToNewHeads(ctx context.Context,
	channel chan<- eth.BlockHeader) (eth.Subscription, error) {
	ch := make(chan *types.Header)
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
	return c.b.SubscribeNewHead(context.Background(), ch)
}
