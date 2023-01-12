// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package access_controller

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
)

var AccessControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060fa8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80636b14daf814602d575b600080fd5b604160383660046055565b60019392505050565b604051901515815260200160405180910390f35b600080600060408486031215606957600080fd5b833573ffffffffffffffffffffffffffffffffffffffff81168114608c57600080fd5b9250602084013567ffffffffffffffff8082111560a857600080fd5b818601915086601f83011260bb57600080fd5b81358181111560c957600080fd5b87602082850101111560da57600080fd5b602083019450809350505050925092509256fea164736f6c6343000806000a",
}

var AccessControllerABI = AccessControllerMetaData.ABI

var AccessControllerBin = AccessControllerMetaData.Bin

func DeployAccessController(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AccessController, error) {
	parsed, err := AccessControllerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AccessControllerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessController{AccessControllerCaller: AccessControllerCaller{contract: contract}, AccessControllerTransactor: AccessControllerTransactor{contract: contract}, AccessControllerFilterer: AccessControllerFilterer{contract: contract}}, nil
}

type AccessController struct {
	address common.Address
	abi     abi.ABI
	AccessControllerCaller
	AccessControllerTransactor
	AccessControllerFilterer
}

type AccessControllerCaller struct {
	contract *bind.BoundContract
}

type AccessControllerTransactor struct {
	contract *bind.BoundContract
}

type AccessControllerFilterer struct {
	contract *bind.BoundContract
}

type AccessControllerSession struct {
	Contract     *AccessController
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AccessControllerCallerSession struct {
	Contract *AccessControllerCaller
	CallOpts bind.CallOpts
}

type AccessControllerTransactorSession struct {
	Contract     *AccessControllerTransactor
	TransactOpts bind.TransactOpts
}

type AccessControllerRaw struct {
	Contract *AccessController
}

type AccessControllerCallerRaw struct {
	Contract *AccessControllerCaller
}

type AccessControllerTransactorRaw struct {
	Contract *AccessControllerTransactor
}

func NewAccessController(address common.Address, backend bind.ContractBackend) (*AccessController, error) {
	abi, err := abi.JSON(strings.NewReader(AccessControllerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAccessController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessController{address: address, abi: abi, AccessControllerCaller: AccessControllerCaller{contract: contract}, AccessControllerTransactor: AccessControllerTransactor{contract: contract}, AccessControllerFilterer: AccessControllerFilterer{contract: contract}}, nil
}

func NewAccessControllerCaller(address common.Address, caller bind.ContractCaller) (*AccessControllerCaller, error) {
	contract, err := bindAccessController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerCaller{contract: contract}, nil
}

func NewAccessControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControllerTransactor, error) {
	contract, err := bindAccessController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControllerTransactor{contract: contract}, nil
}

func NewAccessControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControllerFilterer, error) {
	contract, err := bindAccessController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControllerFilterer{contract: contract}, nil
}

func bindAccessController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControllerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_AccessController *AccessControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessController.Contract.AccessControllerCaller.contract.Call(opts, result, method, params...)
}

func (_AccessController *AccessControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessController.Contract.AccessControllerTransactor.contract.Transfer(opts)
}

func (_AccessController *AccessControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessController.Contract.AccessControllerTransactor.contract.Transact(opts, method, params...)
}

func (_AccessController *AccessControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessController.Contract.contract.Call(opts, result, method, params...)
}

func (_AccessController *AccessControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessController.Contract.contract.Transfer(opts)
}

func (_AccessController *AccessControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessController.Contract.contract.Transact(opts, method, params...)
}

func (_AccessController *AccessControllerCaller) HasAccess(opts *bind.CallOpts, arg0 common.Address, arg1 []byte) (bool, error) {
	var out []interface{}
	err := _AccessController.contract.Call(opts, &out, "hasAccess", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AccessController *AccessControllerSession) HasAccess(arg0 common.Address, arg1 []byte) (bool, error) {
	return _AccessController.Contract.HasAccess(&_AccessController.CallOpts, arg0, arg1)
}

func (_AccessController *AccessControllerCallerSession) HasAccess(arg0 common.Address, arg1 []byte) (bool, error) {
	return _AccessController.Contract.HasAccess(&_AccessController.CallOpts, arg0, arg1)
}

func (_AccessController *AccessController) Address() common.Address {
	return _AccessController.address
}

type AccessControllerInterface interface {
	HasAccess(opts *bind.CallOpts, arg0 common.Address, arg1 []byte) (bool, error)

	Address() common.Address
}
