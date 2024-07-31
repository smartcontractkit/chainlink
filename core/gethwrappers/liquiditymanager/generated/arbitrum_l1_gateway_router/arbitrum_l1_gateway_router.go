// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_l1_gateway_router

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

var ArbitrumL1GatewayRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1ERC20\",\"type\":\"address\"}],\"name\":\"calculateL2TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"finalizeInboundTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"getOutboundCalldata\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inbox\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"outboundTransfer\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_refundTo\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"outboundTransferCustomRefund\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_gateway\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxSubmissionCost\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_creditBackAddress\",\"type\":\"address\"}],\"name\":\"setGateway\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_gateway\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxSubmissionCost\",\"type\":\"uint256\"}],\"name\":\"setGateway\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var ArbitrumL1GatewayRouterABI = ArbitrumL1GatewayRouterMetaData.ABI

type ArbitrumL1GatewayRouter struct {
	address common.Address
	abi     abi.ABI
	ArbitrumL1GatewayRouterCaller
	ArbitrumL1GatewayRouterTransactor
	ArbitrumL1GatewayRouterFilterer
}

type ArbitrumL1GatewayRouterCaller struct {
	contract *bind.BoundContract
}

type ArbitrumL1GatewayRouterTransactor struct {
	contract *bind.BoundContract
}

type ArbitrumL1GatewayRouterFilterer struct {
	contract *bind.BoundContract
}

type ArbitrumL1GatewayRouterSession struct {
	Contract     *ArbitrumL1GatewayRouter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbitrumL1GatewayRouterCallerSession struct {
	Contract *ArbitrumL1GatewayRouterCaller
	CallOpts bind.CallOpts
}

type ArbitrumL1GatewayRouterTransactorSession struct {
	Contract     *ArbitrumL1GatewayRouterTransactor
	TransactOpts bind.TransactOpts
}

type ArbitrumL1GatewayRouterRaw struct {
	Contract *ArbitrumL1GatewayRouter
}

type ArbitrumL1GatewayRouterCallerRaw struct {
	Contract *ArbitrumL1GatewayRouterCaller
}

type ArbitrumL1GatewayRouterTransactorRaw struct {
	Contract *ArbitrumL1GatewayRouterTransactor
}

func NewArbitrumL1GatewayRouter(address common.Address, backend bind.ContractBackend) (*ArbitrumL1GatewayRouter, error) {
	abi, err := abi.JSON(strings.NewReader(ArbitrumL1GatewayRouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbitrumL1GatewayRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL1GatewayRouter{address: address, abi: abi, ArbitrumL1GatewayRouterCaller: ArbitrumL1GatewayRouterCaller{contract: contract}, ArbitrumL1GatewayRouterTransactor: ArbitrumL1GatewayRouterTransactor{contract: contract}, ArbitrumL1GatewayRouterFilterer: ArbitrumL1GatewayRouterFilterer{contract: contract}}, nil
}

func NewArbitrumL1GatewayRouterCaller(address common.Address, caller bind.ContractCaller) (*ArbitrumL1GatewayRouterCaller, error) {
	contract, err := bindArbitrumL1GatewayRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL1GatewayRouterCaller{contract: contract}, nil
}

func NewArbitrumL1GatewayRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrumL1GatewayRouterTransactor, error) {
	contract, err := bindArbitrumL1GatewayRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL1GatewayRouterTransactor{contract: contract}, nil
}

func NewArbitrumL1GatewayRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrumL1GatewayRouterFilterer, error) {
	contract, err := bindArbitrumL1GatewayRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrumL1GatewayRouterFilterer{contract: contract}, nil
}

func bindArbitrumL1GatewayRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbitrumL1GatewayRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumL1GatewayRouter.Contract.ArbitrumL1GatewayRouterCaller.contract.Call(opts, result, method, params...)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.ArbitrumL1GatewayRouterTransactor.contract.Transfer(opts)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.ArbitrumL1GatewayRouterTransactor.contract.Transact(opts, method, params...)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumL1GatewayRouter.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.contract.Transfer(opts)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.contract.Transact(opts, method, params...)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCaller) CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumL1GatewayRouter.contract.Call(opts, &out, "calculateL2TokenAddress", l1ERC20)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _ArbitrumL1GatewayRouter.Contract.CalculateL2TokenAddress(&_ArbitrumL1GatewayRouter.CallOpts, l1ERC20)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCallerSession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _ArbitrumL1GatewayRouter.Contract.CalculateL2TokenAddress(&_ArbitrumL1GatewayRouter.CallOpts, l1ERC20)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCaller) GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	var out []interface{}
	err := _ArbitrumL1GatewayRouter.contract.Call(opts, &out, "getOutboundCalldata", _token, _from, _to, _amount, _data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _ArbitrumL1GatewayRouter.Contract.GetOutboundCalldata(&_ArbitrumL1GatewayRouter.CallOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCallerSession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _ArbitrumL1GatewayRouter.Contract.GetOutboundCalldata(&_ArbitrumL1GatewayRouter.CallOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCaller) Inbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumL1GatewayRouter.contract.Call(opts, &out, "inbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) Inbox() (common.Address, error) {
	return _ArbitrumL1GatewayRouter.Contract.Inbox(&_ArbitrumL1GatewayRouter.CallOpts)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCallerSession) Inbox() (common.Address, error) {
	return _ArbitrumL1GatewayRouter.Contract.Inbox(&_ArbitrumL1GatewayRouter.CallOpts)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumL1GatewayRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) Owner() (common.Address, error) {
	return _ArbitrumL1GatewayRouter.Contract.Owner(&_ArbitrumL1GatewayRouter.CallOpts)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCallerSession) Owner() (common.Address, error) {
	return _ArbitrumL1GatewayRouter.Contract.Owner(&_ArbitrumL1GatewayRouter.CallOpts)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ArbitrumL1GatewayRouter.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ArbitrumL1GatewayRouter.Contract.SupportsInterface(&_ArbitrumL1GatewayRouter.CallOpts, interfaceId)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ArbitrumL1GatewayRouter.Contract.SupportsInterface(&_ArbitrumL1GatewayRouter.CallOpts, interfaceId)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactor) FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.contract.Transact(opts, "finalizeInboundTransfer", _token, _from, _to, _amount, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.FinalizeInboundTransfer(&_ArbitrumL1GatewayRouter.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactorSession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.FinalizeInboundTransfer(&_ArbitrumL1GatewayRouter.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactor) OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.contract.Transact(opts, "outboundTransfer", _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.OutboundTransfer(&_ArbitrumL1GatewayRouter.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactorSession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.OutboundTransfer(&_ArbitrumL1GatewayRouter.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactor) OutboundTransferCustomRefund(opts *bind.TransactOpts, _token common.Address, _refundTo common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.contract.Transact(opts, "outboundTransferCustomRefund", _token, _refundTo, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) OutboundTransferCustomRefund(_token common.Address, _refundTo common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.OutboundTransferCustomRefund(&_ArbitrumL1GatewayRouter.TransactOpts, _token, _refundTo, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactorSession) OutboundTransferCustomRefund(_token common.Address, _refundTo common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.OutboundTransferCustomRefund(&_ArbitrumL1GatewayRouter.TransactOpts, _token, _refundTo, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactor) SetGateway(opts *bind.TransactOpts, _gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int, _creditBackAddress common.Address) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.contract.Transact(opts, "setGateway", _gateway, _maxGas, _gasPriceBid, _maxSubmissionCost, _creditBackAddress)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) SetGateway(_gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int, _creditBackAddress common.Address) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.SetGateway(&_ArbitrumL1GatewayRouter.TransactOpts, _gateway, _maxGas, _gasPriceBid, _maxSubmissionCost, _creditBackAddress)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactorSession) SetGateway(_gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int, _creditBackAddress common.Address) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.SetGateway(&_ArbitrumL1GatewayRouter.TransactOpts, _gateway, _maxGas, _gasPriceBid, _maxSubmissionCost, _creditBackAddress)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactor) SetGateway0(opts *bind.TransactOpts, _gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.contract.Transact(opts, "setGateway0", _gateway, _maxGas, _gasPriceBid, _maxSubmissionCost)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterSession) SetGateway0(_gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.SetGateway0(&_ArbitrumL1GatewayRouter.TransactOpts, _gateway, _maxGas, _gasPriceBid, _maxSubmissionCost)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouterTransactorSession) SetGateway0(_gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int) (*types.Transaction, error) {
	return _ArbitrumL1GatewayRouter.Contract.SetGateway0(&_ArbitrumL1GatewayRouter.TransactOpts, _gateway, _maxGas, _gasPriceBid, _maxSubmissionCost)
}

func (_ArbitrumL1GatewayRouter *ArbitrumL1GatewayRouter) Address() common.Address {
	return _ArbitrumL1GatewayRouter.address
}

type ArbitrumL1GatewayRouterInterface interface {
	CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error)

	GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error)

	Inbox(opts *bind.CallOpts) (common.Address, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error)

	OutboundTransferCustomRefund(opts *bind.TransactOpts, _token common.Address, _refundTo common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error)

	SetGateway(opts *bind.TransactOpts, _gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int, _creditBackAddress common.Address) (*types.Transaction, error)

	SetGateway0(opts *bind.TransactOpts, _gateway common.Address, _maxGas *big.Int, _gasPriceBid *big.Int, _maxSubmissionCost *big.Int) (*types.Transaction, error)

	Address() common.Address
}
