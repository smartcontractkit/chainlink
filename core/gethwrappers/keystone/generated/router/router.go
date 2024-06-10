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
	Bin: "0x608060405234801561001057600080fd5b5033806000816100675760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615610097576100978161009f565b505050610148565b336001600160a01b038216036100f75760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161005e565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b610aa6806101576000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806379ba509711610076578063abcef5541161005b578063abcef554146101c7578063e6b7145814610200578063f2fde38b1461023657600080fd5b806379ba5097146101805780638da5cb5b1461018857600080fd5b80634d93172d116100a75780634d93172d14610138578063516db4081461014d5780635c41d2fe1461016d57600080fd5b8063181f5a77146100c3578063233fd52d14610115575b600080fd5b6100ff6040518060400160405280601481526020017f4b657973746f6e65526f7574657220312e302e3000000000000000000000000081525081565b60405161010c9190610834565b60405180910390f35b610128610123366004610912565b610249565b604051901515815260200161010c565b61014b6101463660046109ad565b610441565b005b61016061015b3660046109cf565b6104bd565b60405161010c91906109e8565b61014b61017b3660046109ad565b61052c565b61014b6105ab565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b6101286101d53660046109ad565b73ffffffffffffffffffffffffffffffffffffffff1660009081526002602052604090205460ff1690565b6101a261020e3660046109cf565b60009081526003602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b61014b6102443660046109ad565b6106a8565b3360009081526002602052604081205460ff16610292576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008881526003602052604090205473ffffffffffffffffffffffffffffffffffffffff16156102f6576040517fa53dc8ca000000000000000000000000000000000000000000000000000000008152600481018990526024015b60405180910390fd5b600088815260036020526040812080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8a81169190911790915587163b900361035857506000610436565b6040517f805f213200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87169063805f2132906103b0908890889088908890600401610a72565b600060405180830381600087803b1580156103ca57600080fd5b505af19250505080156103db575060015b6103e757506000610436565b50600087815260036020526040902080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000017905560015b979650505050505050565b6104496106bc565b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fb96d15bf9258c7b8df062753a6a262864611fc7b060a5ee2e57e79b85f898d389190a250565b60008181526003602052604081205473ffffffffffffffffffffffffffffffffffffffff166104ee57506000919050565b60008281526003602052604090205474010000000000000000000000000000000000000000900460ff16610523576002610526565b60015b92915050565b6105346106bc565b73ffffffffffffffffffffffffffffffffffffffff811660008181526002602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0ea0ce2c048ff45a4a95f2947879de3fb94abec2f152190400cab2d1272a68e79190a250565b60015473ffffffffffffffffffffffffffffffffffffffff16331461062c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102ed565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6106b06106bc565b6106b98161073f565b50565b60005473ffffffffffffffffffffffffffffffffffffffff16331461073d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102ed565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036107be576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102ed565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b600060208083528351808285015260005b8181101561086157858101830151858201604001528201610845565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b803573ffffffffffffffffffffffffffffffffffffffff811681146108c457600080fd5b919050565b60008083601f8401126108db57600080fd5b50813567ffffffffffffffff8111156108f357600080fd5b60208301915083602082850101111561090b57600080fd5b9250929050565b600080600080600080600060a0888a03121561092d57600080fd5b8735965061093d602089016108a0565b955061094b604089016108a0565b9450606088013567ffffffffffffffff8082111561096857600080fd5b6109748b838c016108c9565b909650945060808a013591508082111561098d57600080fd5b5061099a8a828b016108c9565b989b979a50959850939692959293505050565b6000602082840312156109bf57600080fd5b6109c8826108a0565b9392505050565b6000602082840312156109e157600080fd5b5035919050565b6020810160038310610a23577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b604081526000610a86604083018688610a29565b8281036020840152610436818587610a2956fea164736f6c6343000813000a",
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
