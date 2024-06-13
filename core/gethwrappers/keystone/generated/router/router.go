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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"AlreadyAttempted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"ForwarderRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"addForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"getTransmissionState\",\"outputs\":[{\"internalType\":\"enumIRouter.TransmissionState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"}],\"name\":\"getTransmitter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"isForwarder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"forwarder\",\"type\":\"address\"}],\"name\":\"removeForwarder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"transmissionId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"route\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b5033805f816100655760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b5f80546001600160a01b0319166001600160a01b0384811691909117909155811615610094576100948161009c565b505050610144565b336001600160a01b038216036100f45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005c565b600180546001600160a01b0319166001600160a01b038381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610a75806101515f395ff3fe608060405234801561000f575f80fd5b50600436106100b9575f3560e01c806379ba509711610072578063abcef55411610058578063abcef554146101c0578063e6b71458146101f8578063f2fde38b1461022d575f80fd5b806379ba50971461017a5780638da5cb5b14610182575f80fd5b80634d93172d116100a25780634d93172d14610132578063516db408146101475780635c41d2fe14610167575f80fd5b8063181f5a77146100bd578063233fd52d1461010f575b5f80fd5b6100f96040518060400160405280601481526020017f4b657973746f6e65526f7574657220312e302e3000000000000000000000000081525081565b604051610106919061081a565b60405180910390f35b61012261011d3660046108f1565b610240565b6040519015158152602001610106565b610145610140366004610985565b61042f565b005b61015a6101553660046109a5565b6104aa565b60405161010691906109bc565b610145610175366004610985565b610516565b610145610594565b5f5473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610106565b6101226101ce366004610985565b73ffffffffffffffffffffffffffffffffffffffff165f9081526002602052604090205460ff1690565b61019b6102063660046109a5565b5f9081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b61014561023b366004610985565b610690565b335f9081526002602052604081205460ff16610288576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8881526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16156102eb576040517fa53dc8ca000000000000000000000000000000000000000000000000000000008152600481018990526024015b60405180910390fd5b5f88815260036020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8a81169190911790915587163b900361034b57505f610424565b6040517f805f213200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87169063805f2132906103a3908890889088908890600401610a42565b5f604051808303815f87803b1580156103ba575f80fd5b505af19250505080156103cb575060015b6103d657505f610424565b505f87815260036020526040902080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000017905560015b979650505050505050565b6104376106a4565b73ffffffffffffffffffffffffffffffffffffffff81165f8181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d389190a250565b5f8181526003602052604081205473ffffffffffffffffffffffffffffffffffffffff166104d957505f919050565b5f8281526003602052604090205474010000000000000000000000000000000000000000900460ff1661050d576002610510565b60015b92915050565b61051e6106a4565b73ffffffffffffffffffffffffffffffffffffffff81165f8181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e79190a250565b60015473ffffffffffffffffffffffffffffffffffffffff163314610615576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102e2565b5f8054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6106986106a4565b6106a181610726565b50565b5f5473ffffffffffffffffffffffffffffffffffffffff163314610724576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102e2565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036107a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102e2565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8381169182179092555f8054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b5f602080835283518060208501525f5b818110156108465785810183015185820160400152820161082a565b505f6040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff811681146108a7575f80fd5b919050565b5f8083601f8401126108bc575f80fd5b50813567ffffffffffffffff8111156108d3575f80fd5b6020830191508360208285010111156108ea575f80fd5b9250929050565b5f805f805f805f60a0888a031215610907575f80fd5b8735965061091760208901610884565b955061092560408901610884565b9450606088013567ffffffffffffffff80821115610941575f80fd5b61094d8b838c016108ac565b909650945060808a0135915080821115610965575f80fd5b506109728a828b016108ac565b989b979a50959850939692959293505050565b5f60208284031215610995575f80fd5b61099e82610884565b9392505050565b5f602082840312156109b5575f80fd5b5035919050565b60208101600383106109f5577f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b91905290565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b604081525f610a556040830186886109fb565b82810360208401526104248185876109fb56fea164736f6c6343000818000a",
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

func (_KeystoneRouter *KeystoneRouterTransactor) Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte) (*types.Transaction, error) {
	return _KeystoneRouter.contract.Transact(opts, "route", transmissionId, transmitter, receiver, metadata, report)
}

func (_KeystoneRouter *KeystoneRouterSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.Route(&_KeystoneRouter.TransactOpts, transmissionId, transmitter, receiver, metadata, report)
}

func (_KeystoneRouter *KeystoneRouterTransactorSession) Route(transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte) (*types.Transaction, error) {
	return _KeystoneRouter.Contract.Route(&_KeystoneRouter.TransactOpts, transmissionId, transmitter, receiver, metadata, report)
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

	Route(opts *bind.TransactOpts, transmissionId [32]byte, transmitter common.Address, receiver common.Address, metadata []byte, report []byte) (*types.Transaction, error)

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
