// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abstract_arbitrum_token_gateway

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

var AbstractArbitrumTokenGatewayMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1ERC20\",\"type\":\"address\"}],\"name\":\"calculateL2TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counterpartGateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"finalizeInboundTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"getOutboundCalldata\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"outboundTransfer\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var AbstractArbitrumTokenGatewayABI = AbstractArbitrumTokenGatewayMetaData.ABI

type AbstractArbitrumTokenGateway struct {
	address common.Address
	abi     abi.ABI
	AbstractArbitrumTokenGatewayCaller
	AbstractArbitrumTokenGatewayTransactor
	AbstractArbitrumTokenGatewayFilterer
}

type AbstractArbitrumTokenGatewayCaller struct {
	contract *bind.BoundContract
}

type AbstractArbitrumTokenGatewayTransactor struct {
	contract *bind.BoundContract
}

type AbstractArbitrumTokenGatewayFilterer struct {
	contract *bind.BoundContract
}

type AbstractArbitrumTokenGatewaySession struct {
	Contract     *AbstractArbitrumTokenGateway
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AbstractArbitrumTokenGatewayCallerSession struct {
	Contract *AbstractArbitrumTokenGatewayCaller
	CallOpts bind.CallOpts
}

type AbstractArbitrumTokenGatewayTransactorSession struct {
	Contract     *AbstractArbitrumTokenGatewayTransactor
	TransactOpts bind.TransactOpts
}

type AbstractArbitrumTokenGatewayRaw struct {
	Contract *AbstractArbitrumTokenGateway
}

type AbstractArbitrumTokenGatewayCallerRaw struct {
	Contract *AbstractArbitrumTokenGatewayCaller
}

type AbstractArbitrumTokenGatewayTransactorRaw struct {
	Contract *AbstractArbitrumTokenGatewayTransactor
}

func NewAbstractArbitrumTokenGateway(address common.Address, backend bind.ContractBackend) (*AbstractArbitrumTokenGateway, error) {
	abi, err := abi.JSON(strings.NewReader(AbstractArbitrumTokenGatewayABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAbstractArbitrumTokenGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AbstractArbitrumTokenGateway{address: address, abi: abi, AbstractArbitrumTokenGatewayCaller: AbstractArbitrumTokenGatewayCaller{contract: contract}, AbstractArbitrumTokenGatewayTransactor: AbstractArbitrumTokenGatewayTransactor{contract: contract}, AbstractArbitrumTokenGatewayFilterer: AbstractArbitrumTokenGatewayFilterer{contract: contract}}, nil
}

func NewAbstractArbitrumTokenGatewayCaller(address common.Address, caller bind.ContractCaller) (*AbstractArbitrumTokenGatewayCaller, error) {
	contract, err := bindAbstractArbitrumTokenGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AbstractArbitrumTokenGatewayCaller{contract: contract}, nil
}

func NewAbstractArbitrumTokenGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*AbstractArbitrumTokenGatewayTransactor, error) {
	contract, err := bindAbstractArbitrumTokenGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AbstractArbitrumTokenGatewayTransactor{contract: contract}, nil
}

func NewAbstractArbitrumTokenGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*AbstractArbitrumTokenGatewayFilterer, error) {
	contract, err := bindAbstractArbitrumTokenGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AbstractArbitrumTokenGatewayFilterer{contract: contract}, nil
}

func bindAbstractArbitrumTokenGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AbstractArbitrumTokenGatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AbstractArbitrumTokenGateway.Contract.AbstractArbitrumTokenGatewayCaller.contract.Call(opts, result, method, params...)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.AbstractArbitrumTokenGatewayTransactor.contract.Transfer(opts)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.AbstractArbitrumTokenGatewayTransactor.contract.Transact(opts, method, params...)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AbstractArbitrumTokenGateway.Contract.contract.Call(opts, result, method, params...)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.contract.Transfer(opts)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.contract.Transact(opts, method, params...)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCaller) CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error) {
	var out []interface{}
	err := _AbstractArbitrumTokenGateway.contract.Call(opts, &out, "calculateL2TokenAddress", l1ERC20)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewaySession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _AbstractArbitrumTokenGateway.Contract.CalculateL2TokenAddress(&_AbstractArbitrumTokenGateway.CallOpts, l1ERC20)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCallerSession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _AbstractArbitrumTokenGateway.Contract.CalculateL2TokenAddress(&_AbstractArbitrumTokenGateway.CallOpts, l1ERC20)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCaller) CounterpartGateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AbstractArbitrumTokenGateway.contract.Call(opts, &out, "counterpartGateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewaySession) CounterpartGateway() (common.Address, error) {
	return _AbstractArbitrumTokenGateway.Contract.CounterpartGateway(&_AbstractArbitrumTokenGateway.CallOpts)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCallerSession) CounterpartGateway() (common.Address, error) {
	return _AbstractArbitrumTokenGateway.Contract.CounterpartGateway(&_AbstractArbitrumTokenGateway.CallOpts)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCaller) GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	var out []interface{}
	err := _AbstractArbitrumTokenGateway.contract.Call(opts, &out, "getOutboundCalldata", _token, _from, _to, _amount, _data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewaySession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _AbstractArbitrumTokenGateway.Contract.GetOutboundCalldata(&_AbstractArbitrumTokenGateway.CallOpts, _token, _from, _to, _amount, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCallerSession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _AbstractArbitrumTokenGateway.Contract.GetOutboundCalldata(&_AbstractArbitrumTokenGateway.CallOpts, _token, _from, _to, _amount, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCaller) Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AbstractArbitrumTokenGateway.contract.Call(opts, &out, "router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewaySession) Router() (common.Address, error) {
	return _AbstractArbitrumTokenGateway.Contract.Router(&_AbstractArbitrumTokenGateway.CallOpts)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayCallerSession) Router() (common.Address, error) {
	return _AbstractArbitrumTokenGateway.Contract.Router(&_AbstractArbitrumTokenGateway.CallOpts)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayTransactor) FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.contract.Transact(opts, "finalizeInboundTransfer", _token, _from, _to, _amount, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewaySession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.FinalizeInboundTransfer(&_AbstractArbitrumTokenGateway.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayTransactorSession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.FinalizeInboundTransfer(&_AbstractArbitrumTokenGateway.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayTransactor) OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.contract.Transact(opts, "outboundTransfer", _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewaySession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.OutboundTransfer(&_AbstractArbitrumTokenGateway.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGatewayTransactorSession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _AbstractArbitrumTokenGateway.Contract.OutboundTransfer(&_AbstractArbitrumTokenGateway.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_AbstractArbitrumTokenGateway *AbstractArbitrumTokenGateway) Address() common.Address {
	return _AbstractArbitrumTokenGateway.address
}

type AbstractArbitrumTokenGatewayInterface interface {
	CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error)

	CounterpartGateway(opts *bind.CallOpts) (common.Address, error)

	GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error)

	Router(opts *bind.CallOpts) (common.Address, error)

	FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error)

	Address() common.Address
}
