// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_dispute_game_factory

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type IOptimismDisputeGameFactoryGameSearchResult struct {
	Index     *big.Int
	Metadata  [32]byte
	Timestamp uint64
	RootClaim [32]byte
	ExtraData []byte
}

var OptimismDisputeGameFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"GameType\",\"name\":\"_gameType\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"_start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_n\",\"type\":\"uint256\"}],\"name\":\"findLatestGames\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"GameId\",\"name\":\"metadata\",\"type\":\"bytes32\"},{\"internalType\":\"Timestamp\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"Claim\",\"name\":\"rootClaim\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structIOptimismDisputeGameFactory.GameSearchResult[]\",\"name\":\"games_\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gameCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"gameCount_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var OptimismDisputeGameFactoryABI = OptimismDisputeGameFactoryMetaData.ABI

type OptimismDisputeGameFactory struct {
	address common.Address
	abi     abi.ABI
	OptimismDisputeGameFactoryCaller
	OptimismDisputeGameFactoryTransactor
	OptimismDisputeGameFactoryFilterer
}

type OptimismDisputeGameFactoryCaller struct {
	contract *bind.BoundContract
}

type OptimismDisputeGameFactoryTransactor struct {
	contract *bind.BoundContract
}

type OptimismDisputeGameFactoryFilterer struct {
	contract *bind.BoundContract
}

type OptimismDisputeGameFactorySession struct {
	Contract     *OptimismDisputeGameFactory
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismDisputeGameFactoryCallerSession struct {
	Contract *OptimismDisputeGameFactoryCaller
	CallOpts bind.CallOpts
}

type OptimismDisputeGameFactoryTransactorSession struct {
	Contract     *OptimismDisputeGameFactoryTransactor
	TransactOpts bind.TransactOpts
}

type OptimismDisputeGameFactoryRaw struct {
	Contract *OptimismDisputeGameFactory
}

type OptimismDisputeGameFactoryCallerRaw struct {
	Contract *OptimismDisputeGameFactoryCaller
}

type OptimismDisputeGameFactoryTransactorRaw struct {
	Contract *OptimismDisputeGameFactoryTransactor
}

func NewOptimismDisputeGameFactory(address common.Address, backend bind.ContractBackend) (*OptimismDisputeGameFactory, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismDisputeGameFactoryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismDisputeGameFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismDisputeGameFactory{address: address, abi: abi, OptimismDisputeGameFactoryCaller: OptimismDisputeGameFactoryCaller{contract: contract}, OptimismDisputeGameFactoryTransactor: OptimismDisputeGameFactoryTransactor{contract: contract}, OptimismDisputeGameFactoryFilterer: OptimismDisputeGameFactoryFilterer{contract: contract}}, nil
}

func NewOptimismDisputeGameFactoryCaller(address common.Address, caller bind.ContractCaller) (*OptimismDisputeGameFactoryCaller, error) {
	contract, err := bindOptimismDisputeGameFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismDisputeGameFactoryCaller{contract: contract}, nil
}

func NewOptimismDisputeGameFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismDisputeGameFactoryTransactor, error) {
	contract, err := bindOptimismDisputeGameFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismDisputeGameFactoryTransactor{contract: contract}, nil
}

func NewOptimismDisputeGameFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismDisputeGameFactoryFilterer, error) {
	contract, err := bindOptimismDisputeGameFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismDisputeGameFactoryFilterer{contract: contract}, nil
}

func bindOptimismDisputeGameFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismDisputeGameFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismDisputeGameFactory.Contract.OptimismDisputeGameFactoryCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismDisputeGameFactory.Contract.OptimismDisputeGameFactoryTransactor.contract.Transfer(opts)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismDisputeGameFactory.Contract.OptimismDisputeGameFactoryTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismDisputeGameFactory.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismDisputeGameFactory.Contract.contract.Transfer(opts)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismDisputeGameFactory.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryCaller) FindLatestGames(opts *bind.CallOpts, _gameType uint32, _start *big.Int, _n *big.Int) ([]IOptimismDisputeGameFactoryGameSearchResult, error) {
	var out []interface{}
	err := _OptimismDisputeGameFactory.contract.Call(opts, &out, "findLatestGames", _gameType, _start, _n)

	if err != nil {
		return *new([]IOptimismDisputeGameFactoryGameSearchResult), err
	}

	out0 := *abi.ConvertType(out[0], new([]IOptimismDisputeGameFactoryGameSearchResult)).(*[]IOptimismDisputeGameFactoryGameSearchResult)

	return out0, err

}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactorySession) FindLatestGames(_gameType uint32, _start *big.Int, _n *big.Int) ([]IOptimismDisputeGameFactoryGameSearchResult, error) {
	return _OptimismDisputeGameFactory.Contract.FindLatestGames(&_OptimismDisputeGameFactory.CallOpts, _gameType, _start, _n)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryCallerSession) FindLatestGames(_gameType uint32, _start *big.Int, _n *big.Int) ([]IOptimismDisputeGameFactoryGameSearchResult, error) {
	return _OptimismDisputeGameFactory.Contract.FindLatestGames(&_OptimismDisputeGameFactory.CallOpts, _gameType, _start, _n)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryCaller) GameCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OptimismDisputeGameFactory.contract.Call(opts, &out, "gameCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactorySession) GameCount() (*big.Int, error) {
	return _OptimismDisputeGameFactory.Contract.GameCount(&_OptimismDisputeGameFactory.CallOpts)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactoryCallerSession) GameCount() (*big.Int, error) {
	return _OptimismDisputeGameFactory.Contract.GameCount(&_OptimismDisputeGameFactory.CallOpts)
}

func (_OptimismDisputeGameFactory *OptimismDisputeGameFactory) Address() common.Address {
	return _OptimismDisputeGameFactory.address
}

type OptimismDisputeGameFactoryInterface interface {
	FindLatestGames(opts *bind.CallOpts, _gameType uint32, _start *big.Int, _n *big.Int) ([]IOptimismDisputeGameFactoryGameSearchResult, error)

	GameCount(opts *bind.CallOpts) (*big.Int, error)

	Address() common.Address
}
