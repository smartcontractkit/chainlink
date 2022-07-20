// Package blockchain handles connections to various blockchains
package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EVMClient is the interface that wraps a given client implementation for a blockchain, to allow for switching
// of network types within the test suite
// EVMClient can be connected to a single or multiple nodes,
type EVMClient interface {
	// Getters
	Get() interface{}
	GetNetworkName() string
	NetworkSimulated() bool
	GetChainID() *big.Int
	GetClients() []EVMClient
	GetDefaultWallet() *EthereumWallet
	GetWallets() []*EthereumWallet
	GetNetworkConfig() *EVMNetwork

	// Setters
	SetID(id int)
	SetDefaultWallet(num int) error
	SetWallets([]*EthereumWallet)
	LoadWallets(ns interface{}) error
	SwitchNode(node int) error

	// On-chain Operations
	BalanceAt(ctx context.Context, address common.Address) (*big.Int, error)
	HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error)
	HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error)
	LatestBlockNumber(ctx context.Context) (uint64, error)
	Fund(toAddress string, amount *big.Float) error
	DeployContract(
		contractName string,
		deployer ContractDeployer,
	) (*common.Address, *types.Transaction, interface{}, error)
	TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error)
	ProcessTransaction(tx *types.Transaction) error
	IsTxConfirmed(txHash common.Hash) (bool, error)
	GetTxReceipt(txHash common.Hash) (*types.Receipt, error)
	ParallelTransactions(enabled bool)
	Close() error

	// Gas Operations
	EstimateCostForChainlinkOperations(amountOfOperations int) (*big.Float, error)
	EstimateTransactionGasCost() (*big.Int, error)
	GasStats() *GasStats

	// Event Subscriptions
	AddHeaderEventSubscription(key string, subscriber HeaderEventSubscription)
	DeleteHeaderEventSubscription(key string)
	WaitForEvents() error
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

// NodeBlock block with a node ID which mined it
type NodeBlock struct {
	NodeID int
	*types.Block
}

// HeaderEventSubscription is an interface for allowing callbacks when the client receives a new header
type HeaderEventSubscription interface {
	ReceiveBlock(header NodeBlock) error
	Wait() error
}

// ContractDeployer acts as a go-between function for general contract deployment
type ContractDeployer func(auth *bind.TransactOpts, backend bind.ContractBackend) (
	common.Address,
	*types.Transaction,
	interface{},
	error,
)
