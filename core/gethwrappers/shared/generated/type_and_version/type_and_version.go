// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package type_and_version

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

var ITypeAndVersionMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"}]",
}

var ITypeAndVersionABI = ITypeAndVersionMetaData.ABI

type ITypeAndVersion struct {
	address common.Address
	abi     abi.ABI
	ITypeAndVersionCaller
	ITypeAndVersionTransactor
	ITypeAndVersionFilterer
}

type ITypeAndVersionCaller struct {
	contract *bind.BoundContract
}

type ITypeAndVersionTransactor struct {
	contract *bind.BoundContract
}

type ITypeAndVersionFilterer struct {
	contract *bind.BoundContract
}

type ITypeAndVersionSession struct {
	Contract     *ITypeAndVersion
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ITypeAndVersionCallerSession struct {
	Contract *ITypeAndVersionCaller
	CallOpts bind.CallOpts
}

type ITypeAndVersionTransactorSession struct {
	Contract     *ITypeAndVersionTransactor
	TransactOpts bind.TransactOpts
}

type ITypeAndVersionRaw struct {
	Contract *ITypeAndVersion
}

type ITypeAndVersionCallerRaw struct {
	Contract *ITypeAndVersionCaller
}

type ITypeAndVersionTransactorRaw struct {
	Contract *ITypeAndVersionTransactor
}

func NewITypeAndVersion(address common.Address, backend bind.ContractBackend) (*ITypeAndVersion, error) {
	abi, err := abi.JSON(strings.NewReader(ITypeAndVersionABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindITypeAndVersion(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ITypeAndVersion{address: address, abi: abi, ITypeAndVersionCaller: ITypeAndVersionCaller{contract: contract}, ITypeAndVersionTransactor: ITypeAndVersionTransactor{contract: contract}, ITypeAndVersionFilterer: ITypeAndVersionFilterer{contract: contract}}, nil
}

func NewITypeAndVersionCaller(address common.Address, caller bind.ContractCaller) (*ITypeAndVersionCaller, error) {
	contract, err := bindITypeAndVersion(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ITypeAndVersionCaller{contract: contract}, nil
}

func NewITypeAndVersionTransactor(address common.Address, transactor bind.ContractTransactor) (*ITypeAndVersionTransactor, error) {
	contract, err := bindITypeAndVersion(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ITypeAndVersionTransactor{contract: contract}, nil
}

func NewITypeAndVersionFilterer(address common.Address, filterer bind.ContractFilterer) (*ITypeAndVersionFilterer, error) {
	contract, err := bindITypeAndVersion(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ITypeAndVersionFilterer{contract: contract}, nil
}

func bindITypeAndVersion(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ITypeAndVersionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ITypeAndVersion *ITypeAndVersionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ITypeAndVersion.Contract.ITypeAndVersionCaller.contract.Call(opts, result, method, params...)
}

func (_ITypeAndVersion *ITypeAndVersionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITypeAndVersion.Contract.ITypeAndVersionTransactor.contract.Transfer(opts)
}

func (_ITypeAndVersion *ITypeAndVersionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITypeAndVersion.Contract.ITypeAndVersionTransactor.contract.Transact(opts, method, params...)
}

func (_ITypeAndVersion *ITypeAndVersionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ITypeAndVersion.Contract.contract.Call(opts, result, method, params...)
}

func (_ITypeAndVersion *ITypeAndVersionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITypeAndVersion.Contract.contract.Transfer(opts)
}

func (_ITypeAndVersion *ITypeAndVersionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITypeAndVersion.Contract.contract.Transact(opts, method, params...)
}

func (_ITypeAndVersion *ITypeAndVersionCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ITypeAndVersion.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_ITypeAndVersion *ITypeAndVersionSession) TypeAndVersion() (string, error) {
	return _ITypeAndVersion.Contract.TypeAndVersion(&_ITypeAndVersion.CallOpts)
}

func (_ITypeAndVersion *ITypeAndVersionCallerSession) TypeAndVersion() (string, error) {
	return _ITypeAndVersion.Contract.TypeAndVersion(&_ITypeAndVersion.CallOpts)
}

func (_ITypeAndVersion *ITypeAndVersion) Address() common.Address {
	return _ITypeAndVersion.address
}

type ITypeAndVersionInterface interface {
	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Address() common.Address
}
