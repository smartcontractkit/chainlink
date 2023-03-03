// Package blockchain handles connections to various blockchains
package blockchain

import (
	"context"
	"crypto/ecdsa"
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
	GetNonceSetting() NonceSettings

	// Setters
	SetID(id int)
	SetDefaultWallet(num int) error
	SetWallets([]*EthereumWallet)
	LoadWallets(ns EVMNetwork) error
	SwitchNode(node int) error
	SyncNonce(c EVMClient)

	// On-chain Operations
	BalanceAt(ctx context.Context, address common.Address) (*big.Int, error)
	HeaderHashByNumber(ctx context.Context, bn *big.Int) (string, error)
	HeaderTimestampByNumber(ctx context.Context, bn *big.Int) (uint64, error)
	LatestBlockNumber(ctx context.Context) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	Fund(toAddress string, amount *big.Float) error
	ReturnFunds(fromKey *ecdsa.PrivateKey) error
	DeployContract(
		contractName string,
		deployer ContractDeployer,
	) (*common.Address, *types.Transaction, interface{}, error)
	LoadContract(contractName string, address common.Address, loader ContractLoader) (interface{}, error)
	TransactionOpts(from *EthereumWallet) (*bind.TransactOpts, error)
	ProcessTransaction(tx *types.Transaction) error
	ProcessEvent(name string, event *types.Log, confirmedChan chan bool, errorChan chan error) error
	IsEventConfirmed(event *types.Log) (confirmed, removed bool, err error)
	IsTxConfirmed(txHash common.Hash) (bool, error)
	GetTxReceipt(txHash common.Hash) (*types.Receipt, error)
	ParallelTransactions(enabled bool)
	Close() error
	Backend() bind.ContractBackend

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

// NodeHeader header with the ID of the node that received it
type NodeHeader struct {
	NodeID int
	types.Header
}

// HeaderEventSubscription is an interface for allowing callbacks when the client receives a new header
type HeaderEventSubscription interface {
	ReceiveHeader(header NodeHeader) error
	Wait() error
	Complete() bool
}

// ContractDeployer acts as a go-between function for general contract deployment
type ContractDeployer func(auth *bind.TransactOpts, backend bind.ContractBackend) (
	common.Address,
	*types.Transaction,
	interface{},
	error,
)

// ContractLoader helps loading already deployed contracts
type ContractLoader func(address common.Address, backend bind.ContractBackend) (
	interface{},
	error,
)
