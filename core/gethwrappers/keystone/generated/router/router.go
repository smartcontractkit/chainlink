// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package router

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var KeystoneRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"AlreadyAttempted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"addForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"getTransmissionState\",\"outputs\":[{\"internalType\":\"enumIRouter.TransmissionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"removeForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"gasForCallExactCheck\",\"type\":\"uint16\"}],\"name\":\"route\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610bae806101576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638da5cb5b11610076578063bddd690d1161005b578063bddd690d146101ed578063e6b7145814610200578063f2fde38b1461023657600080fd5b80638da5cb5b14610165578063abcef554146101a457600080fd5b8063516db408116100a7578063516db4081461012a5780635c41d2fe1461014a57806379ba50971461015d57600080fd5b8063181f5a77146100c35780634d93172d14610115575b600080fd5b6100ff6040518060400160405280601481526020017f4b657973746f6e65526f7574657220312e302e3000000000000000000000000081525081565b60405161010c919061090e565b60405180910390f35b6101286101233660046109a3565b610249565b005b61013d6101383660046109c5565b6102c5565b60405161010c91906109de565b6101286101583660046109a3565b610334565b6101286103b3565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b6101dd6101b23660046109a3565b73ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b604051901515815260200161010c565b6101dd6101fb366004610a68565b6104b5565b61017f61020e3660046109c5565b60009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b6101286102443660046109a3565b6106cc565b6102516106e0565b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d389190a250565b60008181526003602052604081205473ffffffffffffffffffffffffffffffffffffffff166102f657506000919050565b60008281526003602052604090205474010000000000000000000000000000000000000000900460ff1661032b57600261032e565b60015b92915050565b61033c6106e0565b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e79190a250565b60015473ffffffffffffffffffffffffffffffffffffffff163314610439576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b3360009081526002602052604081205460ff166104fe576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008a81526003602052604090205473ffffffffffffffffffffffffffffffffffffffff161561055d576040517fa53dc8ca000000000000000000000000000000000000000000000000000000008152600481018b9052602401610430565b60008a815260036020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8c81169190911790915589163b90036105bf575060006106bf565b600061066863805f213260e01b898989896040516024016105e39493929190610b6f565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091528a8686610763565b905080156106bc5760008b815260036020526040902080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790555b90505b9998505050505050505050565b6106d46106e0565b6106dd81610819565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610761576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610430565b565b6000833b610795577f0c3b563c0000000000000000000000000000000000000000000000000000000060005260046000fd5b5a828110156107c8577fafa32a2c0000000000000000000000000000000000000000000000000000000060005260046000fd5b8290036040810481038410610801577f37c3be290000000000000000000000000000000000000000000000000000000060005260046000fd5b5060008086516020880160008888f195945050505050565b3373ffffffffffffffffffffffffffffffffffffffff821603610898576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610430565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208083528351808285015260005b8181101561093b5785810183015185820160400152820161091f565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461099e57600080fd5b919050565b6000602082840312156109b557600080fd5b6109be8261097a565b9392505050565b6000602082840312156109d757600080fd5b5035919050565b6020810160038310610a19577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b60008083601f840112610a3157600080fd5b50813567ffffffffffffffff811115610a4957600080fd5b602083019150836020828501011115610a6157600080fd5b9250929050565b600080600080600080600080600060e08a8c031215610a8657600080fd5b89359850610a9660208b0161097a565b9750610aa460408b0161097a565b965060608a013567ffffffffffffffff80821115610ac157600080fd5b610acd8d838e01610a1f565b909850965060808c0135915080821115610ae657600080fd5b50610af38c828d01610a1f565b90955093505060a08a0135915060c08a013561ffff81168114610b1557600080fd5b809150509295985092959850929598565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b604081526000610b83604083018688610b26565b8281036020840152610b96818587610b26565b97965050505050505056fea164736f6c6343000813000a",
}

var KeystoneRouterABI = KeystoneRouterMetaData.ABI

var KeystoneRouterBin = KeystoneRouterMetaData.Bin

func DeployKeystoneRouter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KeystoneRouter, error) {
	parsed, err := KeystoneRouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeystoneRouterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KeystoneRouter{address: address, abi: *parsed, KeystoneRouterCaller: KeystoneRouterCaller{contract: contract}, KeystoneRouterTransactor: KeystoneRouterTransactor{contract: contract}, KeystoneRouterFilterer: KeystoneRouterFilterer{contract: contract}}, nil
}

type KeystoneRouter struct {
	address common.Address
	abi     abi.ABI
	KeystoneRouterCaller
	KeystoneRouterTransactor
	KeystoneRouterFilterer
}

type KeystoneRouterCaller struct {
	contract *bind.BoundContract
}

type KeystoneRouterTransactor struct {
	contract *bind.BoundContract
}

type KeystoneRouterFilterer struct {
	contract *bind.BoundContract
}

type KeystoneRouterSession struct {
	Contract     *KeystoneRouter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type KeystoneRouterCallerSession struct {
	Contract *KeystoneRouterCaller
	CallOpts bind.CallOpts
}

type KeystoneRouterTransactorSession struct {
	Contract     *KeystoneRouterTransactor
	TransactOpts bind.TransactOpts
}

type KeystoneRouterRaw struct {
	Contract *KeystoneRouter
}

type KeystoneRouterCallerRaw struct {
	Contract *KeystoneRouterCaller
}

type KeystoneRouterTransactorRaw struct {
	Contract *KeystoneRouterTransactor
}

func NewKeystoneRouter(address common.Address, backend bind.ContractBackend) (*KeystoneRouter, error) {
	abi, err := abi.JSON(strings.NewReader(KeystoneRouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindKeystoneRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouter{address: address, abi: abi, KeystoneRouterCaller: KeystoneRouterCaller{contract: contract}, KeystoneRouterTransactor: KeystoneRouterTransactor{contract: contract}, KeystoneRouterFilterer: KeystoneRouterFilterer{contract: contract}}, nil
}

func NewKeystoneRouterCaller(address common.Address, caller bind.ContractCaller) (*KeystoneRouterCaller, error) {
	contract, err := bindKeystoneRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouterCaller{contract: contract}, nil
}

func NewKeystoneRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*KeystoneRouterTransactor, error) {
	contract, err := bindKeystoneRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouterTransactor{contract: contract}, nil
}

func NewKeystoneRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*KeystoneRouterFilterer, error) {
	contract, err := bindKeystoneRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouterFilterer{contract: contract}, nil
}

func bindKeystoneRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KeystoneRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_KeystoneRouter *KeystoneRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneRouter.Contract.KeystoneRouterCaller.contract.Call(opts, result, method, params...)
}

func (_KeystoneRouter *KeystoneRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.KeystoneRouterTransactor.contract.Transfer(opts)
}

func (_KeystoneRouter *KeystoneRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.KeystoneRouterTransactor.contract.Transact(opts, method, params...)
}

func (_KeystoneRouter *KeystoneRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KeystoneRouter.Contract.contract.Call(opts, result, method, params...)
}

func (_KeystoneRouter *KeystoneRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.contract.Transfer(opts)
}

func (_KeystoneRouter *KeystoneRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.contract.Transact(opts, method, params...)
}

func (_KeystoneRouter *KeystoneRouterCaller) GetTransmissionState(opts *bind.CallOpts, transmissionId [32]byte) (uint8, error) {
	var out []interface{}
	err := _KeystoneRouter.contract.Call(opts, &out, "getTransmissionState", transmissionId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_KeystoneRouter *KeystoneRouterSession) GetTransmissionState(transmissionId [32]byte) (uint8, error) {
	return _KeystoneRouter.Contract.GetTransmissionState(&_KeystoneRouter.CallOpts, transmissionId)
}

func (_KeystoneRouter *KeystoneRouterCallerSession) GetTransmissionState(transmissionId [32]byte) (uint8, error) {
	return _KeystoneRouter.Contract.GetTransmissionState(&_KeystoneRouter.CallOpts, transmissionId)
}

func (_KeystoneRouter *KeystoneRouterCaller) GetTransmitter(opts *bind.CallOpts, transmissionId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _KeystoneRouter.contract.Call(opts, &out, "getTransmitter", transmissionId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneRouter *KeystoneRouterSession) GetTransmitter(transmissionId [32]byte) (common.Address, error) {
	return _KeystoneRouter.Contract.GetTransmitter(&_KeystoneRouter.CallOpts, transmissionId)
}

func (_KeystoneRouter *KeystoneRouterCallerSession) GetTransmitter(transmissionId [32]byte) (common.Address, error) {
	return _KeystoneRouter.Contract.GetTransmitter(&_KeystoneRouter.CallOpts, transmissionId)
}

func (_KeystoneRouter *KeystoneRouterCaller) IsForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error) {
	var out []interface{}
	err := _KeystoneRouter.contract.Call(opts, &out, "isForwarder", forwarder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_KeystoneRouter *KeystoneRouterSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _KeystoneRouter.Contract.IsForwarder(&_KeystoneRouter.CallOpts, forwarder)
}

func (_KeystoneRouter *KeystoneRouterCallerSession) IsForwarder(forwarder common.Address) (bool, error) {
	return _KeystoneRouter.Contract.IsForwarder(&_KeystoneRouter.CallOpts, forwarder)
}

func (_KeystoneRouter *KeystoneRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KeystoneRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_KeystoneRouter *KeystoneRouterSession) Owner() (common.Address, error) {
	return _KeystoneRouter.Contract.Owner(&_KeystoneRouter.CallOpts)
}

func (_KeystoneRouter *KeystoneRouterCallerSession) Owner() (common.Address, error) {
	return _KeystoneRouter.Contract.Owner(&_KeystoneRouter.CallOpts)
}

func (_KeystoneRouter *KeystoneRouterCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KeystoneRouter.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_KeystoneRouter *KeystoneRouterSession) TypeAndVersion() (string, error) {
	return _KeystoneRouter.Contract.TypeAndVersion(&_KeystoneRouter.CallOpts)
}

func (_KeystoneRouter *KeystoneRouterCallerSession) TypeAndVersion() (string, error) {
	return _KeystoneRouter.Contract.TypeAndVersion(&_KeystoneRouter.CallOpts)
}

func (_KeystoneRouter *KeystoneRouterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeystoneRouter.contract.Transact(opts, "acceptOwnership")
}

func (_KeystoneRouter *KeystoneRouterSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneRouter.Contract.AcceptOwnership(&_KeystoneRouter.TransactOpts)
}

func (_KeystoneRouter *KeystoneRouterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _KeystoneRouter.Contract.AcceptOwnership(&_KeystoneRouter.TransactOpts)
}

func (_KeystoneRouter *KeystoneRouterTransactor) AddForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.contract.Transact(opts, "addForwarder", forwarder)
}

func (_KeystoneRouter *KeystoneRouterSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.AddForwarder(&_KeystoneRouter.TransactOpts, forwarder)
}

func (_KeystoneRouter *KeystoneRouterTransactorSession) AddForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.AddForwarder(&_KeystoneRouter.TransactOpts, forwarder)
}

func (_KeystoneRouter *KeystoneRouterTransactor) RemoveForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.contract.Transact(opts, "removeForwarder", forwarder)
}

func (_KeystoneRouter *KeystoneRouterSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.RemoveForwarder(&_KeystoneRouter.TransactOpts, forwarder)
}

func (_KeystoneRouter *KeystoneRouterTransactorSession) RemoveForwarder(forwarder common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.RemoveForwarder(&_KeystoneRouter.TransactOpts, forwarder)
}

func (_KeystoneRouter *KeystoneRouterTransactor) Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte, gasLimit *big.Int, gasForCallExactCheck uint16) (*types.Transaction, error) {
	return _KeystoneRouter.contract.Transact(opts, "route", transmissionId, transmitter, receiver, metadata, report, gasLimit, gasForCallExactCheck)
}

func (_KeystoneRouter *KeystoneRouterSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte, gasLimit *big.Int, gasForCallExactCheck uint16) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.Route(&_KeystoneRouter.TransactOpts, transmissionId, transmitter, receiver, metadata, report, gasLimit, gasForCallExactCheck)
}

func (_KeystoneRouter *KeystoneRouterTransactorSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte, gasLimit *big.Int, gasForCallExactCheck uint16) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.Route(&_KeystoneRouter.TransactOpts, transmissionId, transmitter, receiver, metadata, report, gasLimit, gasForCallExactCheck)
}

func (_KeystoneRouter *KeystoneRouterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.contract.Transact(opts, "transferOwnership", to)
}

func (_KeystoneRouter *KeystoneRouterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.TransferOwnership(&_KeystoneRouter.TransactOpts, to)
}

func (_KeystoneRouter *KeystoneRouterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.TransferOwnership(&_KeystoneRouter.TransactOpts, to)
}

type KeystoneRouterForwarderAddedIterator struct {
	Event *KeystoneRouterForwarderAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneRouterForwarderAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneRouterForwarderAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeystoneRouterForwarderAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeystoneRouterForwarderAddedIterator) Error() error {
	return it.fail
}

func (it *KeystoneRouterForwarderAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneRouterForwarderAdded struct {
	Forwarder common.Address
	Raw       types.Log
}

func (_KeystoneRouter *KeystoneRouterFilterer) FilterForwarderAdded(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneRouterForwarderAddedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneRouter.contract.FilterLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouterForwarderAddedIterator{contract: _KeystoneRouter.contract, event: "ForwarderAdded", logs: logs, sub: sub}, nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) WatchForwarderAdded(opts *bind.WatchOpts, sink chan<- *KeystoneRouterForwarderAdded, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneRouter.contract.WatchLogs(opts, "ForwarderAdded", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneRouterForwarderAdded)
				if err := _KeystoneRouter.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) ParseForwarderAdded(log types.Log) (*KeystoneRouterForwarderAdded, error) {
	event := new(KeystoneRouterForwarderAdded)
	if err := _KeystoneRouter.contract.UnpackLog(event, "ForwarderAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneRouterForwarderRemovedIterator struct {
	Event *KeystoneRouterForwarderRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneRouterForwarderRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneRouterForwarderRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeystoneRouterForwarderRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeystoneRouterForwarderRemovedIterator) Error() error {
	return it.fail
}

func (it *KeystoneRouterForwarderRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneRouterForwarderRemoved struct {
	Forwarder common.Address
	Raw       types.Log
}

func (_KeystoneRouter *KeystoneRouterFilterer) FilterForwarderRemoved(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneRouterForwarderRemovedIterator, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneRouter.contract.FilterLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouterForwarderRemovedIterator{contract: _KeystoneRouter.contract, event: "ForwarderRemoved", logs: logs, sub: sub}, nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) WatchForwarderRemoved(opts *bind.WatchOpts, sink chan<- *KeystoneRouterForwarderRemoved, forwarder []common.Address) (event.Subscription, error) {

	var forwarderRule []interface{}
	for _, forwarderItem := range forwarder {
		forwarderRule = append(forwarderRule, forwarderItem)
	}

	logs, sub, err := _KeystoneRouter.contract.WatchLogs(opts, "ForwarderRemoved", forwarderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneRouterForwarderRemoved)
				if err := _KeystoneRouter.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) ParseForwarderRemoved(log types.Log) (*KeystoneRouterForwarderRemoved, error) {
	event := new(KeystoneRouterForwarderRemoved)
	if err := _KeystoneRouter.contract.UnpackLog(event, "ForwarderRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneRouterOwnershipTransferRequestedIterator struct {
	Event *KeystoneRouterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneRouterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneRouterOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeystoneRouterOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeystoneRouterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *KeystoneRouterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneRouterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneRouter *KeystoneRouterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneRouterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneRouter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouterOwnershipTransferRequestedIterator{contract: _KeystoneRouter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneRouter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneRouterOwnershipTransferRequested)
				if err := _KeystoneRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) ParseOwnershipTransferRequested(log types.Log) (*KeystoneRouterOwnershipTransferRequested, error) {
	event := new(KeystoneRouterOwnershipTransferRequested)
	if err := _KeystoneRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type KeystoneRouterOwnershipTransferredIterator struct {
	Event *KeystoneRouterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *KeystoneRouterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeystoneRouterOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(KeystoneRouterOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *KeystoneRouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *KeystoneRouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type KeystoneRouterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_KeystoneRouter *KeystoneRouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneRouterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneRouter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &KeystoneRouterOwnershipTransferredIterator{contract: _KeystoneRouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _KeystoneRouter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(KeystoneRouterOwnershipTransferred)
				if err := _KeystoneRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_KeystoneRouter *KeystoneRouterFilterer) ParseOwnershipTransferred(log types.Log) (*KeystoneRouterOwnershipTransferred, error) {
	event := new(KeystoneRouterOwnershipTransferred)
	if err := _KeystoneRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_KeystoneRouter *KeystoneRouter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _KeystoneRouter.abi.Events["ForwarderAdded"].ID:
		return _KeystoneRouter.ParseForwarderAdded(log)
	case _KeystoneRouter.abi.Events["ForwarderRemoved"].ID:
		return _KeystoneRouter.ParseForwarderRemoved(log)
	case _KeystoneRouter.abi.Events["OwnershipTransferRequested"].ID:
		return _KeystoneRouter.ParseOwnershipTransferRequested(log)
	case _KeystoneRouter.abi.Events["OwnershipTransferred"].ID:
		return _KeystoneRouter.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (KeystoneRouterForwarderAdded) Topic() common.Hash {
	return common.HexToHash("0x0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e7")
}

func (KeystoneRouterForwarderRemoved) Topic() common.Hash {
	return common.HexToHash("0xb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d38")
}

func (KeystoneRouterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (KeystoneRouterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_KeystoneRouter *KeystoneRouter) Address() common.Address {
	return _KeystoneRouter.address
}

type KeystoneRouterInterface interface {
	GetTransmissionState(opts *bind.CallOpts, transmissionId [32]byte) (uint8, error)

	GetTransmitter(opts *bind.CallOpts, transmissionId [32]byte) (common.Address, error)

	IsForwarder(opts *bind.CallOpts, forwarder common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error)

	RemoveForwarder(opts *bind.TransactOpts, forwarder common.Address) (*types.Transaction, error)

	Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte, gasLimit *big.Int, gasForCallExactCheck uint16) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterForwarderAdded(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneRouterForwarderAddedIterator, error)

	WatchForwarderAdded(opts *bind.WatchOpts, sink chan<- *KeystoneRouterForwarderAdded, forwarder []common.Address) (event.Subscription, error)

	ParseForwarderAdded(log types.Log) (*KeystoneRouterForwarderAdded, error)

	FilterForwarderRemoved(opts *bind.FilterOpts, forwarder []common.Address) (*KeystoneRouterForwarderRemovedIterator, error)

	WatchForwarderRemoved(opts *bind.WatchOpts, sink chan<- *KeystoneRouterForwarderRemoved, forwarder []common.Address) (event.Subscription, error)

	ParseForwarderRemoved(log types.Log) (*KeystoneRouterForwarderRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneRouterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *KeystoneRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*KeystoneRouterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*KeystoneRouterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *KeystoneRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*KeystoneRouterOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
