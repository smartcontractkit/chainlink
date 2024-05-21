// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrfv2plus_client

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

var VRFV2PlusClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"EXTRA_ARGS_V1_TAG\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a0610038600b82828239805160001a607314602b57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361060335760003560e01c8063f7514ab4146038575b600080fd5b605e7f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa81565b6040517fffffffff00000000000000000000000000000000000000000000000000000000909116815260200160405180910390f3fea164736f6c6343000813000a",
}

var VRFV2PlusClientABI = VRFV2PlusClientMetaData.ABI

var VRFV2PlusClientBin = VRFV2PlusClientMetaData.Bin

func DeployVRFV2PlusClient(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFV2PlusClient, error) {
	parsed, err := VRFV2PlusClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFV2PlusClientBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFV2PlusClient{address: address, abi: *parsed, VRFV2PlusClientCaller: VRFV2PlusClientCaller{contract: contract}, VRFV2PlusClientTransactor: VRFV2PlusClientTransactor{contract: contract}, VRFV2PlusClientFilterer: VRFV2PlusClientFilterer{contract: contract}}, nil
}

type VRFV2PlusClient struct {
	address common.Address
	abi     abi.ABI
	VRFV2PlusClientCaller
	VRFV2PlusClientTransactor
	VRFV2PlusClientFilterer
}

type VRFV2PlusClientCaller struct {
	contract *bind.BoundContract
}

type VRFV2PlusClientTransactor struct {
	contract *bind.BoundContract
}

type VRFV2PlusClientFilterer struct {
	contract *bind.BoundContract
}

type VRFV2PlusClientSession struct {
	Contract     *VRFV2PlusClient
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFV2PlusClientCallerSession struct {
	Contract *VRFV2PlusClientCaller
	CallOpts bind.CallOpts
}

type VRFV2PlusClientTransactorSession struct {
	Contract     *VRFV2PlusClientTransactor
	TransactOpts bind.TransactOpts
}

type VRFV2PlusClientRaw struct {
	Contract *VRFV2PlusClient
}

type VRFV2PlusClientCallerRaw struct {
	Contract *VRFV2PlusClientCaller
}

type VRFV2PlusClientTransactorRaw struct {
	Contract *VRFV2PlusClientTransactor
}

func NewVRFV2PlusClient(address common.Address, backend bind.ContractBackend) (*VRFV2PlusClient, error) {
	abi, err := abi.JSON(strings.NewReader(VRFV2PlusClientABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFV2PlusClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusClient{address: address, abi: abi, VRFV2PlusClientCaller: VRFV2PlusClientCaller{contract: contract}, VRFV2PlusClientTransactor: VRFV2PlusClientTransactor{contract: contract}, VRFV2PlusClientFilterer: VRFV2PlusClientFilterer{contract: contract}}, nil
}

func NewVRFV2PlusClientCaller(address common.Address, caller bind.ContractCaller) (*VRFV2PlusClientCaller, error) {
	contract, err := bindVRFV2PlusClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusClientCaller{contract: contract}, nil
}

func NewVRFV2PlusClientTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFV2PlusClientTransactor, error) {
	contract, err := bindVRFV2PlusClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusClientTransactor{contract: contract}, nil
}

func NewVRFV2PlusClientFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFV2PlusClientFilterer, error) {
	contract, err := bindVRFV2PlusClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFV2PlusClientFilterer{contract: contract}, nil
}

func bindVRFV2PlusClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFV2PlusClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFV2PlusClient *VRFV2PlusClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusClient.Contract.VRFV2PlusClientCaller.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusClient *VRFV2PlusClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusClient.Contract.VRFV2PlusClientTransactor.contract.Transfer(opts)
}

func (_VRFV2PlusClient *VRFV2PlusClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusClient.Contract.VRFV2PlusClientTransactor.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusClient *VRFV2PlusClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFV2PlusClient.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFV2PlusClient *VRFV2PlusClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFV2PlusClient.Contract.contract.Transfer(opts)
}

func (_VRFV2PlusClient *VRFV2PlusClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFV2PlusClient.Contract.contract.Transact(opts, method, params...)
}

func (_VRFV2PlusClient *VRFV2PlusClientCaller) EXTRAARGSV1TAG(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _VRFV2PlusClient.contract.Call(opts, &out, "EXTRA_ARGS_V1_TAG")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

func (_VRFV2PlusClient *VRFV2PlusClientSession) EXTRAARGSV1TAG() ([4]byte, error) {
	return _VRFV2PlusClient.Contract.EXTRAARGSV1TAG(&_VRFV2PlusClient.CallOpts)
}

func (_VRFV2PlusClient *VRFV2PlusClientCallerSession) EXTRAARGSV1TAG() ([4]byte, error) {
	return _VRFV2PlusClient.Contract.EXTRAARGSV1TAG(&_VRFV2PlusClient.CallOpts)
}

func (_VRFV2PlusClient *VRFV2PlusClient) Address() common.Address {
	return _VRFV2PlusClient.address
}

type VRFV2PlusClientInterface interface {
	EXTRAARGSV1TAG(opts *bind.CallOpts) ([4]byte, error)

	Address() common.Address
}
