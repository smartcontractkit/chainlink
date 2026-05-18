// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_log_automation

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

type Log struct {
	Index       *big.Int
	Timestamp   *big.Int
	TxHash      [32]byte
	BlockNumber *big.Int
	BlockHash   [32]byte
	Source      common.Address
	Topics      [][32]byte
	Data        []byte
}

var ILogAutomationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"topics\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLog\",\"name\":\"log\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"checkData\",\"type\":\"bytes\"}],\"name\":\"checkLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var ILogAutomationABI = ILogAutomationMetaData.ABI

type ILogAutomation struct {
	address common.Address
	abi     abi.ABI
	ILogAutomationCaller
	ILogAutomationTransactor
	ILogAutomationFilterer
}

type ILogAutomationCaller struct {
	contract *bind.BoundContract
}

type ILogAutomationTransactor struct {
	contract *bind.BoundContract
}

type ILogAutomationFilterer struct {
	contract *bind.BoundContract
}

type ILogAutomationSession struct {
	Contract     *ILogAutomation
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ILogAutomationCallerSession struct {
	Contract *ILogAutomationCaller
	CallOpts bind.CallOpts
}

type ILogAutomationTransactorSession struct {
	Contract     *ILogAutomationTransactor
	TransactOpts bind.TransactOpts
}

type ILogAutomationRaw struct {
	Contract *ILogAutomation
}

type ILogAutomationCallerRaw struct {
	Contract *ILogAutomationCaller
}

type ILogAutomationTransactorRaw struct {
	Contract *ILogAutomationTransactor
}

func NewILogAutomation(address common.Address, backend bind.ContractBackend) (*ILogAutomation, error) {
	abi, err := abi.JSON(strings.NewReader(ILogAutomationABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindILogAutomation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ILogAutomation{address: address, abi: abi, ILogAutomationCaller: ILogAutomationCaller{contract: contract}, ILogAutomationTransactor: ILogAutomationTransactor{contract: contract}, ILogAutomationFilterer: ILogAutomationFilterer{contract: contract}}, nil
}

func NewILogAutomationCaller(address common.Address, caller bind.ContractCaller) (*ILogAutomationCaller, error) {
	contract, err := bindILogAutomation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ILogAutomationCaller{contract: contract}, nil
}

func NewILogAutomationTransactor(address common.Address, transactor bind.ContractTransactor) (*ILogAutomationTransactor, error) {
	contract, err := bindILogAutomation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ILogAutomationTransactor{contract: contract}, nil
}

func NewILogAutomationFilterer(address common.Address, filterer bind.ContractFilterer) (*ILogAutomationFilterer, error) {
	contract, err := bindILogAutomation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ILogAutomationFilterer{contract: contract}, nil
}

func bindILogAutomation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ILogAutomationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ILogAutomation *ILogAutomationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILogAutomation.Contract.ILogAutomationCaller.contract.Call(opts, result, method, params...)
}

func (_ILogAutomation *ILogAutomationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILogAutomation.Contract.ILogAutomationTransactor.contract.Transfer(opts)
}

func (_ILogAutomation *ILogAutomationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILogAutomation.Contract.ILogAutomationTransactor.contract.Transact(opts, method, params...)
}

func (_ILogAutomation *ILogAutomationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILogAutomation.Contract.contract.Call(opts, result, method, params...)
}

func (_ILogAutomation *ILogAutomationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILogAutomation.Contract.contract.Transfer(opts)
}

func (_ILogAutomation *ILogAutomationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILogAutomation.Contract.contract.Transact(opts, method, params...)
}

func (_ILogAutomation *ILogAutomationTransactor) CheckLog(opts *bind.TransactOpts, log Log, checkData []byte) (*types.Transaction, error) {
	return _ILogAutomation.contract.Transact(opts, "checkLog", log, checkData)
}

func (_ILogAutomation *ILogAutomationSession) CheckLog(log Log, checkData []byte) (*types.Transaction, error) {
	return _ILogAutomation.Contract.CheckLog(&_ILogAutomation.TransactOpts, log, checkData)
}

func (_ILogAutomation *ILogAutomationTransactorSession) CheckLog(log Log, checkData []byte) (*types.Transaction, error) {
	return _ILogAutomation.Contract.CheckLog(&_ILogAutomation.TransactOpts, log, checkData)
}

func (_ILogAutomation *ILogAutomationTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _ILogAutomation.contract.Transact(opts, "performUpkeep", performData)
}

func (_ILogAutomation *ILogAutomationSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _ILogAutomation.Contract.PerformUpkeep(&_ILogAutomation.TransactOpts, performData)
}

func (_ILogAutomation *ILogAutomationTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _ILogAutomation.Contract.PerformUpkeep(&_ILogAutomation.TransactOpts, performData)
}

func (_ILogAutomation *ILogAutomation) Address() common.Address {
	return _ILogAutomation.address
}

type ILogAutomationInterface interface {
	CheckLog(opts *bind.TransactOpts, log Log, checkData []byte) (*types.Transaction, error)

	PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error)

	Address() common.Address
}
