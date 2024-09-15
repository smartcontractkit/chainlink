// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_token_gateway

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

var ArbitrumTokenGatewayMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1ERC20\",\"type\":\"address\"}],\"name\":\"calculateL2TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"finalizeInboundTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"getOutboundCalldata\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"outboundTransfer\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

var ArbitrumTokenGatewayABI = ArbitrumTokenGatewayMetaData.ABI

type ArbitrumTokenGateway struct {
	address common.Address
	abi     abi.ABI
	ArbitrumTokenGatewayCaller
	ArbitrumTokenGatewayTransactor
	ArbitrumTokenGatewayFilterer
}

type ArbitrumTokenGatewayCaller struct {
	contract *bind.BoundContract
}

type ArbitrumTokenGatewayTransactor struct {
	contract *bind.BoundContract
}

type ArbitrumTokenGatewayFilterer struct {
	contract *bind.BoundContract
}

type ArbitrumTokenGatewaySession struct {
	Contract     *ArbitrumTokenGateway
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbitrumTokenGatewayCallerSession struct {
	Contract *ArbitrumTokenGatewayCaller
	CallOpts bind.CallOpts
}

type ArbitrumTokenGatewayTransactorSession struct {
	Contract     *ArbitrumTokenGatewayTransactor
	TransactOpts bind.TransactOpts
}

type ArbitrumTokenGatewayRaw struct {
	Contract *ArbitrumTokenGateway
}

type ArbitrumTokenGatewayCallerRaw struct {
	Contract *ArbitrumTokenGatewayCaller
}

type ArbitrumTokenGatewayTransactorRaw struct {
	Contract *ArbitrumTokenGatewayTransactor
}

func NewArbitrumTokenGateway(address common.Address, backend bind.ContractBackend) (*ArbitrumTokenGateway, error) {
	abi, err := abi.JSON(strings.NewReader(ArbitrumTokenGatewayABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbitrumTokenGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrumTokenGateway{address: address, abi: abi, ArbitrumTokenGatewayCaller: ArbitrumTokenGatewayCaller{contract: contract}, ArbitrumTokenGatewayTransactor: ArbitrumTokenGatewayTransactor{contract: contract}, ArbitrumTokenGatewayFilterer: ArbitrumTokenGatewayFilterer{contract: contract}}, nil
}

func NewArbitrumTokenGatewayCaller(address common.Address, caller bind.ContractCaller) (*ArbitrumTokenGatewayCaller, error) {
	contract, err := bindArbitrumTokenGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumTokenGatewayCaller{contract: contract}, nil
}

func NewArbitrumTokenGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrumTokenGatewayTransactor, error) {
	contract, err := bindArbitrumTokenGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumTokenGatewayTransactor{contract: contract}, nil
}

func NewArbitrumTokenGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrumTokenGatewayFilterer, error) {
	contract, err := bindArbitrumTokenGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrumTokenGatewayFilterer{contract: contract}, nil
}

func bindArbitrumTokenGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbitrumTokenGatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumTokenGateway.Contract.ArbitrumTokenGatewayCaller.contract.Call(opts, result, method, params...)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.ArbitrumTokenGatewayTransactor.contract.Transfer(opts)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.ArbitrumTokenGatewayTransactor.contract.Transact(opts, method, params...)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumTokenGateway.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.contract.Transfer(opts)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.contract.Transact(opts, method, params...)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayCaller) CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumTokenGateway.contract.Call(opts, &out, "calculateL2TokenAddress", l1ERC20)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewaySession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _ArbitrumTokenGateway.Contract.CalculateL2TokenAddress(&_ArbitrumTokenGateway.CallOpts, l1ERC20)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayCallerSession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _ArbitrumTokenGateway.Contract.CalculateL2TokenAddress(&_ArbitrumTokenGateway.CallOpts, l1ERC20)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayCaller) GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	var out []interface{}
	err := _ArbitrumTokenGateway.contract.Call(opts, &out, "getOutboundCalldata", _token, _from, _to, _amount, _data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewaySession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _ArbitrumTokenGateway.Contract.GetOutboundCalldata(&_ArbitrumTokenGateway.CallOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayCallerSession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _ArbitrumTokenGateway.Contract.GetOutboundCalldata(&_ArbitrumTokenGateway.CallOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayTransactor) FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.contract.Transact(opts, "finalizeInboundTransfer", _token, _from, _to, _amount, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewaySession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.FinalizeInboundTransfer(&_ArbitrumTokenGateway.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayTransactorSession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.FinalizeInboundTransfer(&_ArbitrumTokenGateway.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayTransactor) OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.contract.Transact(opts, "outboundTransfer", _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewaySession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.OutboundTransfer(&_ArbitrumTokenGateway.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGatewayTransactorSession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumTokenGateway.Contract.OutboundTransfer(&_ArbitrumTokenGateway.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumTokenGateway *ArbitrumTokenGateway) Address() common.Address {
	return _ArbitrumTokenGateway.address
}

type ArbitrumTokenGatewayInterface interface {
	CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error)

	GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error)

	FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error)

	Address() common.Address
}
